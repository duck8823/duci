package duci

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/pkg/errors"
	"time"
)

type duci struct {
	executor.Executor
	jobService job_service.Service
	github     github.GitHub
}

// New returns duci instance
func New() (*duci, error) {
	jobService, err := job_service.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	github, err := github.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	duci := &duci{
		jobService: jobService,
		github:     github,
	}
	duci.Executor = executor.DefaultExecutorBuilder().
		StartFunc(duci.Start).
		EndFunc(duci.End).
		LogFunc(duci.AppendLog).
		Build()

	return duci, nil
}

// Start represents a function of start job
func (d *duci) Start(ctx context.Context) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	if err := d.jobService.Start(buildJob.ID); err != nil {
		_ = d.jobService.Append(buildJob.ID, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
		return
	}
	_ = d.github.CreateCommitStatus(ctx, github.CommitStatus{
		TargetSource: buildJob.TargetSource,
		State:        github.PENDING,
		Description:  "pending",
		Context:      buildJob.TaskName,
		TargetURL:    buildJob.TargetURL,
	})
}

// AppendLog is a function that print and store log
func (d *duci) AppendLog(ctx context.Context, log job.Log) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
		println(line.Message)
		_ = d.jobService.Append(buildJob.ID, *line)
	}
}

// End represents a function
func (d *duci) End(ctx context.Context, e error) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	if err := d.jobService.Finish(buildJob.ID); err != nil {
		_ = d.jobService.Append(buildJob.ID, job.LogLine{Timestamp: time.Now(), Message: err.Error()})
		return
	}

	switch e {
	case nil:
		_ = d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.SUCCESS,
			Description:  "success",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		})
	case runner.ErrFailure:
		_ = d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.FAILURE,
			Description:  "failure",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		})
	default:
		_ = d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.ERROR,
			Description:  github.Description(fmt.Sprintf("error: %s", e.Error())),
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		})
	}
}
