package router

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/logstore"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

func New() (http.Handler, error) {
	logstoreService, githubService, err := createCommonServices()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner, err := createRunner(logstoreService, githubService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	webhooksCtrl := &controller.JobController{Runner: dockerRunner, GitHub: githubService}
	logCtrl := &controller.LogController{LogStore: logstoreService}

	rtr := chi.NewRouter()
	rtr.Post("/", webhooksCtrl.ServeHTTP)
	rtr.Get("/logs/{uuid}", logCtrl.ServeHTTP)

	return rtr, nil
}

func createCommonServices() (logstore.Service, github.Service, error) {
	logstoreService, err := logstore.New()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	githubService, err := github.New()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return logstoreService, githubService, nil
}

func createRunner(logstoreService logstore.Service, githubService github.Service) (runner.Runner, error) {
	gitClient, err := git.New(application.Config.GitHub.SSHKeyPath)
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
		LogStore:    logstoreService,
	}

	return dockerRunner, nil
}
