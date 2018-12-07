package router

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	health_controller "github.com/duck8823/duci/presentation/controller/health"
	job_controller "github.com/duck8823/duci/presentation/controller/job"
	webhook_controller "github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

// New returns handler of application.
func New() (http.Handler, error) {
	// FIXME: where initialize??
	if err := git.InitializeWithHTTP(func(ctx context.Context, log job.Log) {
		for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
			println(line.Message)
		}
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := github.Initialize(os.Getenv("GITHUB_TOKEN")); err != nil {
		return nil, errors.WithStack(err)
	}

	gh, err := github.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jobService, err := job_service.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	webhook := &webhook_controller.Handler{
		Executor: executor.DefaultExecutorBuilder().
			StartFunc(func(ctx context.Context) {
				buildJob, err := application.BuildJobFromContext(ctx)
				id := job.ID(buildJob.ID)
				if err != nil {
					_ = jobService.Append(id, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
					return
				}
				if err := jobService.Start(id); err != nil {
					_ = jobService.Append(id, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
					return
				}
				_ = gh.CreateCommitStatus(ctx, github.CommitStatus{
					TargetSource: buildJob.TargetSource,
					State:        github.PENDING,
					Description:  "pending",
					Context:      buildJob.TaskName,
					TargetURL:    buildJob.TargetURL,
				})
			}).
			LogFunc(func(ctx context.Context, log job.Log) {
				buildJob, err := application.BuildJobFromContext(ctx)
				id := job.ID(buildJob.ID)
				if err != nil {
					_ = jobService.Append(id, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
					return
				}
				for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
					println(line.Message)
					_ = jobService.Append(id, *line)
				}
			}).
			EndFunc(func(ctx context.Context, e error) {
				buildJob, err := application.BuildJobFromContext(ctx)
				id := job.ID(buildJob.ID)
				if err != nil {
					_ = jobService.Append(id, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
					return
				}
				if err := jobService.Finish(id); err != nil {
					println(fmt.Sprintf("%+v", err))
					return
				}

				switch e {
				case nil:
					_ = gh.CreateCommitStatus(ctx, github.CommitStatus{
						TargetSource: buildJob.TargetSource,
						State:        github.SUCCESS,
						Description:  "success",
						Context:      buildJob.TaskName,
						TargetURL:    buildJob.TargetURL,
					})
				case runner.ErrFailure:
					_ = gh.CreateCommitStatus(ctx, github.CommitStatus{
						TargetSource: buildJob.TargetSource,
						State:        github.FAILURE,
						Description:  "failure",
						Context:      buildJob.TaskName,
						TargetURL:    buildJob.TargetURL,
					})
				default:
					_ = gh.CreateCommitStatus(ctx, github.CommitStatus{
						TargetSource: buildJob.TargetSource,
						State:        github.ERROR,
						Description:  github.Description(fmt.Sprintf("error: %s", e.Error())),
						Context:      buildJob.TaskName,
						TargetURL:    buildJob.TargetURL,
					})
				}
			}).
			Build(),
	}

	job := &job_controller.Handler{
		Service: jobService,
	}

	docker, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	health := &health_controller.Handler{
		Docker: docker,
	}

	rtr := chi.NewRouter()
	rtr.Post("/", webhook.ServeHTTP)
	rtr.Get("/logs/{uuid}", job.ServeHTTP)
	rtr.Get("/health", health.ServeHTTP)

	return rtr, nil
}
