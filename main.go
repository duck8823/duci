package main

import (
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/duck8823/duci/service/github"
	"github.com/duck8823/duci/service/runner"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
)

const AppName = "duci"

func main() {

	githubService, err := github.NewWithEnv()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to create github service.\n%+v", err)
		os.Exit(1)
		return
	}
	dockerClient, err := docker.New()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to create docker client.\n%+v", err)
		os.Exit(1)
		return
	}

	dockerRunner := &runner.DockerRunner{
		Name:        AppName,
		BaseWorkDir: path.Join(os.TempDir(), AppName),
		GitHub:      githubService,
		Docker:      dockerClient,
	}

	ctrl := &controller.JobController{Runner: dockerRunner, GitHub: githubService}

	http.Handle("/", ctrl)

	http.ListenAndServe(":8080", nil)
}
