package main

import (
	"flag"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func main() {
	flag.Var(application.Config, "c", "configuration file path")
	flag.Parse()

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
		Name:        application.Name,
		BaseWorkDir: application.Config.Server.WorkDir,
		GitHub:      githubService,
		Docker:      dockerClient,
	}

	ctrl := &controller.JobController{Runner: dockerRunner, GitHub: githubService}

	http.Handle("/", ctrl)

	if err := http.ListenAndServe(application.Config.Addr(), nil); err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to run server.\n%+v", err)
		os.Exit(1)
		return
	}
}
