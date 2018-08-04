package main

import (
	"flag"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/log"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/git"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func init() {
	flag.Var(application.Config, "c", "configuration file path")
	flag.Parse()

	if err := semaphore.Make(); err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to initialize a semaphore.\n%+v", err)
		os.Exit(1)
		return
	}
}

func main() {
	jobCtrl, err := jobController()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to initialize job controller.\n%+v", err)
		os.Exit(1)
		return
	}

	logCtrl, err := logController()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to initialize log controller.\n%+v", err)
		os.Exit(1)
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/", jobCtrl.ServeHTTP)
	rtr.Get("/logs/{uuid}", logCtrl.ServeHTTP)

	if err := http.ListenAndServe(application.Config.Addr(), rtr); err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to run server.\n%+v", err)
		os.Exit(1)
		return
	}
}

func jobController() (*controller.JobController, error) {
	gitClient, err := git.New(application.Config.Server.SSHKeyPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	githubService, err := github.NewWithEnv()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dockerClient, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner := &runner.DockerRunner{
		Name:        application.Name,
		BaseWorkDir: application.Config.Server.WorkDir,
		Git:         gitClient,
		GitHub:      githubService,
		Docker:      dockerClient,
	}

	return &controller.JobController{Runner: dockerRunner, GitHub: githubService}, nil
}

func logController() (*controller.LogController, error) {
	logService, err := log.NewStoreService()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &controller.LogController{LogService: logService}, nil
}
