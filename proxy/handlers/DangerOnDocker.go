package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"net/http"
	"os"
	"context"
	"github.com/moby/moby/client"
	"github.com/docker/docker/api/types"
)

type DangerOnDocker struct {}

func (d *DangerOnDocker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body."))
		return
	}

	pr := &github.PullRequestEvent{}
	if err := json.Unmarshal(body, pr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read as pull request."))
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create docker client."))
		return
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get container list."))
		return
	}
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID, container.Image)
	}

	_, err = git.PlainClone("/tmp", false, &git.CloneOptions{
		URL:      *pr.Repo.CloneURL,
		Progress: os.Stdout,
	})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to git clone. %s", err)))
		return
	}
}
