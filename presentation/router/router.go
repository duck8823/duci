package router

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/docker"
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/logstore"
	"github.com/duck8823/duci/application/service/runner"
	moby "github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// New returns handler of application.
func New() (http.Handler, error) {
	dockerClient, err := moby.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logstoreService, githubService, err := createCommonServices()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner, err := createRunner(logstoreService, githubService, dockerClient)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	webhooksCtrl := &controller.WebhooksController{Runner: dockerRunner, GitHub: githubService}
	logCtrl := &controller.LogController{LogStore: logstoreService}
	healthCtrl := &controller.HealthController{Docker: docker.New(dockerClient)}

	rtr := chi.NewRouter()
	rtr.Post("/", webhooksCtrl.ServeHTTP)
	rtr.Get("/logs/{uuid}", logCtrl.ServeHTTP)
	rtr.Get("/health", healthCtrl.ServeHTTP)

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

func createRunner(logstoreService logstore.Service, githubService github.Service, dockerClient moby.Client) (runner.Runner, error) {
	gitClient, err := git.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner := &runner.DockerRunner{
		BaseWorkDir: application.Config.Server.WorkDir,
		Git:         gitClient,
		GitHub:      githubService,
		Docker:      dockerClient,
		LogStore:    logstoreService,
	}

	return dockerRunner, nil
}
