package router

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/docker"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/logstore"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/domain/model/job"
	git2 "github.com/duck8823/duci/domain/model/job/target/git"
	github2 "github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/presentation/controller"
	webhook2 "github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

// New returns handler of application.
func New() (http.Handler, error) {
	dockerService, logstoreService, githubService, err := createCommonServices()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner, err := createRunner(logstoreService, githubService, dockerService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	webhooksCtrl := &controller.WebhooksController{Runner: dockerRunner, GitHub: githubService}
	logCtrl := &controller.LogController{LogStore: logstoreService}
	healthCtrl := &controller.HealthController{Docker: dockerService}

	rtr := chi.NewRouter()
	rtr.Post("/", webhooksCtrl.ServeHTTP)
	rtr.Get("/logs/{uuid}", logCtrl.ServeHTTP)
	rtr.Get("/health", healthCtrl.ServeHTTP)

	// using new api
	// FIXME: where initialize??
	if err := git2.InitializeWithHTTP(func(ctx context.Context, log job.Log) {
		for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
			println(line.Message)
		}
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := github2.Initialize(os.Getenv("GITHUB_TOKEN")); err != nil {
		return nil, errors.WithStack(err)
	}

	webhook := &webhook2.Handler{
		Executor: executor.DefaultExecutorBuilder().
			LogFunc(func(ctx context.Context, log job.Log) {
				for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
					println(line.Message)
				}
			}).
			Build(),
	}

	rtr.Post("/webhook", webhook.ServeHTTP)

	return rtr, nil
}

func createCommonServices() (docker.Service, logstore.Service, github.Service, error) {
	dockerService, err := docker.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	logstoreService, err := logstore.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}
	githubService, err := github.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	return dockerService, logstoreService, githubService, nil
}

func createRunner(logstoreService logstore.Service, githubService github.Service, dockerService docker.Service) (runner.Runner, error) {
	gitClient, err := git.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner := &runner.DockerRunner{
		BaseWorkDir: application.Config.Server.WorkDir,
		Git:         gitClient,
		GitHub:      githubService,
		Docker:      dockerService,
		LogStore:    logstoreService,
	}

	return dockerRunner, nil
}
