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
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	logger.Init("minimal_ci", false, false, os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read Payload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pr := &github.IssueCommentEvent{}
		if err := json.Unmarshal(body, pr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Trigger build
		if r.Header.Get("X-GitHub-Event") != "issue_comment" {
			http.Error(w, "payload event type must be issue_comment", http.StatusBadRequest)
			return
		}
		if ! strings.Contains(*pr.Comment.Body, "test") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not build."))
			return
		}

		// Clone git repository
		base := fmt.Sprintf("%v", time.Now().Unix())
		root := fmt.Sprintf("/tmp/%s", base)
		if _, err = git.PlainClone(root, false, &git.CloneOptions{
			URL:      *pr.Repo.CloneURL,
			Progress: os.Stdout,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create tar archive
		tarFile, err := os.OpenFile(root+"/Dockerfile.tar", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writer.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tarFile.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := os.Open(root + "/Dockerfile.tar")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create docker client
		cli, err := client.NewEnvClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build image
		if _, err := cli.ImageBuild(context.Background(), file, types.ImageBuildOptions{
			Tags: []string{base},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create container
		con, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: base,
			Cmd:   []string{"test"},
		}, nil, nil, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Run container
		if err := cli.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = cli.ContainerWait(context.Background(), con.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		out, err := cli.ContainerLogs(context.Background(), con.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove container
		if err := cli.ContainerRemove(context.Background(), con.ID, types.ContainerRemoveOptions{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response console
		buf := new(bytes.Buffer)
		buf.ReadFrom(out)

		respBody, err := json.Marshal(struct {
			Console string`json:"console"`
		}{
			Console: buf.String(),
		})
		if err := cli.ContainerRemove(context.Background(), con.ID, types.ContainerRemoveOptions{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	})

	http.ListenAndServe(":8080", nil)
}
