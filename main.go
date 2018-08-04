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
	gitClient, err := git.New(application.Config.Server.SSHKeyPath)
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to create git client.\n%+v", err)
		os.Exit(1)
		return
	}
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
		Git:         gitClient,
		GitHub:      githubService,
		Docker:      dockerClient,
	}

	jobCtrl := &controller.JobController{Runner: dockerRunner, GitHub: githubService}

	logService, err := log.NewStoreService()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to initialize database.\n%+v", err)
		os.Exit(1)
		return
	}
	defer logService.Close()

	logCtrl := &controller.LogController{LogService: logService}

	rtr := chi.NewRouter()
	rtr.Post("/", jobCtrl.ServeHTTP)
	rtr.Get("/logs/{uuid}", logCtrl.ServeHTTP)

	if err := http.ListenAndServe(application.Config.Addr(), rtr); err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to run server.\n%+v", err)
		os.Exit(1)
		return
	}
}
