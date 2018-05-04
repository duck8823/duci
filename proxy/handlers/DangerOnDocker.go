package handlers

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/go-github/github"
	"github.com/google/logger"
	"github.com/moby/moby/client"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DangerOnDocker struct{}

func (d *DangerOnDocker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body."))
		return
	}

	pr := &github.IssueCommentEvent{}
	if err := json.Unmarshal(body, pr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read as pull request."))
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	base := fmt.Sprintf("%v", time.Now().Unix())
	root := fmt.Sprintf("/tmp/%s", base)
	_, err = git.PlainClone(root, false, &git.CloneOptions{
		URL:      *pr.Repo.CloneURL,
		Progress: os.Stdout,
	})
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	dockerfileContent := `FROM ruby:2.4.1-alpine3.6

RUN apk update && \
    apk upgrade && \
    apk add git

RUN gem install bundler

ADD . .

RUN bundle install

ENTRYPOINT ["bundle", "exec", "danger"]
`

	gemfileContent := `
source "https://rubygems.org"

git_source(:github) {|repo_name| "https://github.com/#{repo_name}" }

gem "danger"
`

	tarFile, err := os.OpenFile(root+"/Dockerfile.tar", os.O_RDWR|os.O_CREATE, 0666)
	defer tarFile.Close()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	writer := tar.NewWriter(tarFile)
	defer writer.Close()

	for _, file := range []struct {
		Name string
		Body string
	}{
		{Name: "Dockerfile", Body: dockerfileContent},
		{Name: "Gemfile", Body: gemfileContent},
	} {
		header := &tar.Header{
			Name: file.Name,
			Mode: 0666,
			Size: int64(len(file.Body)),
		}
		if err := writer.WriteHeader(header); err != nil {
			logger.Fatal(err.Error())
			os.Exit(-1)
		}
		if _, err := writer.Write([]byte(file.Body)); err != nil {
			logger.Fatal(err.Error())
			os.Exit(-1)
		}
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if file, err := os.Open(path); err != nil {
			return err
		} else {
			defer file.Close()

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			header := &tar.Header{
				Name: strings.Replace(file.Name(), "/tmp/"+base, "", -1),
				Mode: 0666,
				Size: info.Size(),
			}
			if err := writer.WriteHeader(header); err != nil {
				logger.Fatal(err.Error())
				os.Exit(-1)
			}
			if _, err := writer.Write(data); err != nil {
				logger.Fatal(err.Error())
				os.Exit(-1)
			}
		}
		return nil
	}); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	if err := writer.Close(); err != nil {
		logger.Fatal(err)
		os.Exit(-1)
	}
	if err := tarFile.Close(); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	file, _ := os.Open(root + "/Dockerfile.tar")
	defer file.Close()

	resp, err := cli.ImageBuild(context.Background(), file, types.ImageBuildOptions{
		Tags: []string{base},
	})
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	} else {
		logger.Info(string(body))
	}

	if pr.Issue.HTMLURL == nil {
		w.WriteHeader(404)
		return
	}

	con, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: base,
		Cmd:   []string{"pr", *pr.Issue.HTMLURL},
		Env:   []string{fmt.Sprintf("DANGER_GITHUB_API_TOKEN=%s", os.Getenv("DANGER_GITHUB_API_TOKEN"))},
	}, nil, nil, "")
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	if err := cli.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{}); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	if _, err = cli.ContainerWait(context.Background(), con.ID); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	out, err := cli.ContainerLogs(context.Background(), con.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}

	io.Copy(os.Stdout, out)

	if err := cli.ContainerRemove(context.Background(), con.ID, types.ContainerRemoveOptions{}); err != nil {
		logger.Fatal(err.Error())
		os.Exit(-1)
	}
}
