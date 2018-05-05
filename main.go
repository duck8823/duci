package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/go-github/github"
	"github.com/google/logger"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	logger.Init("minimal_ci", false, false, os.Stdout)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(tc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read Payload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Error(err, w, nil)
			return
		}

		event := &github.IssueCommentEvent{}
		if err := json.Unmarshal(body, event); err != nil {
			Error(err, w, nil)
			return
		}

		// Trigger build
		githubEvent := r.Header.Get("X-GitHub-Event")
		if githubEvent != "issue_comment" {
			message := fmt.Sprintf("payload event type must be issue_comment. but %s", githubEvent)
			Error(errors.New(message), w, nil)
			return
		}
		if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not build."))
			return
		}
		phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")

		pr, _, err := githubClient.PullRequests.Get(
			context.Background(),
			event.Repo.Owner.GetLogin(),
			event.Repo.GetName(),
			event.Issue.GetNumber(),
		)
		if err != nil {
			Error(err, w, nil)
			return
		}

		// Clone git repository
		base := fmt.Sprintf("%v", time.Now().Unix())
		root := fmt.Sprintf("/tmp/%s", base)
		repo, err := git.PlainClone(root, false, &git.CloneOptions{
			URL:           event.Repo.GetCloneURL(),
			Progress:      os.Stdout,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", pr.Head.GetRef())),
		})
		if err != nil {
			Error(err, w, nil)
			return
		}

		// Pending Status
		ref, err := repo.Head()
		if err != nil {
			Error(err, w, nil)
			return
		}
		statusService := &CommitStatusService{
			Context:      fmt.Sprintf("minimal_ci-%s", phrase),
			GithubClient: githubClient,
			Repo:         event.Repo,
			Hash:         ref.Hash(),
		}
		statusService.Create(PENDING)

		// Create tar archive
		tarFile, err := os.OpenFile(root+"/Dockerfile.tar", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			Error(err, w, statusService)
			return
		}
		defer tarFile.Close()

		writer := tar.NewWriter(tarFile)
		defer writer.Close()

		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			header := &tar.Header{
				Name: strings.Replace(file.Name(), root, "", -1),
				Mode: 0666,
				Size: info.Size(),
			}
			if err := writer.WriteHeader(header); err != nil {
				return err
			}
			if _, err := writer.Write(data); err != nil {
				return err
			}
			return nil
		}); err != nil {
			Error(err, w, statusService)
			return
		}

		if err := writer.Close(); err != nil {
			Error(err, w, statusService)
			return
		}
		if err := tarFile.Close(); err != nil {
			Error(err, w, statusService)
			return
		}

		file, err := os.Open(root + "/Dockerfile.tar")
		if err != nil {
			Error(err, w, statusService)
			return
		}
		defer file.Close()

		// Create docker client
		cli, err := client.NewEnvClient()
		if err != nil {
			Error(err, w, statusService)
			return
		}

		// Build image
		if resp, err := cli.ImageBuild(context.Background(), file, types.ImageBuildOptions{
			Tags: []string{base},
		}); err != nil {
			Error(err, w, statusService)
			return
		} else {
			defer resp.Body.Close()

			if _, err := ioutil.ReadAll(resp.Body); err != nil {
				Error(err, w, statusService)
				return
			}
			logger.Info("Image Build succeeded.")
		}

		// Create container
		con, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: base,
			Cmd:   []string{phrase},
		}, nil, nil, "")
		if err != nil {
			Error(err, w, statusService)
			return
		}

		// Run container
		if err := cli.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{}); err != nil {
			Error(err, w, statusService)
			return
		}

		if code, err := cli.ContainerWait(context.Background(), con.ID); err != nil {
			Error(err, w, statusService)
			return
		} else if code != 0 {
			// Failure status
			statusService.Create(FAILURE)

			http.Error(w, fmt.Sprintf("return code: %v", code), http.StatusInternalServerError)
			return
		}

		out, err := cli.ContainerLogs(context.Background(), con.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			Error(err, w, statusService)
			return
		}

		// Remove container
		if err := cli.ContainerRemove(context.Background(), con.ID, types.ContainerRemoveOptions{}); err != nil {
			Error(err, w, statusService)
			return
		}

		// Remove image
		if _, err := cli.ImageRemove(context.Background(), base, types.ImageRemoveOptions{}); err != nil {
			Error(err, w, statusService)
			return
		}

		// Succeed Status
		statusService.Create(SUCCESS)

		// Response console
		buf := new(bytes.Buffer)
		buf.ReadFrom(out)

		respBody, err := json.Marshal(struct {
			Console string `json:"console"`
		}{
			Console: buf.String(),
		})

		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	})

	http.ListenAndServe(":8080", nil)
}

func Error(err error, w http.ResponseWriter, s *CommitStatusService) {
	logger.Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
	if s != nil {
		s.Create(ERROR)
	}
}

type CommitStatusService struct {
	Context      string
	GithubClient *github.Client
	Repo         *github.Repository
	Hash         plumbing.Hash
}

func (s *CommitStatusService) Create(state State) error {
	str := string(state)
	_, _, err := s.GithubClient.Repositories.CreateStatus(
		context.Background(),
		s.Repo.Owner.GetLogin(),
		s.Repo.GetName(),
		s.Hash.String(),
		&github.RepoStatus{
			Context: &s.Context,
			State:   &str,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

type State string

const (
	PENDING State = "pending"
	SUCCESS State = "success"
	ERROR   State = "error"
	FAILURE State = "failure"
)
