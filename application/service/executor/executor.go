package executor

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
)

// Executor is job executor
type Executor interface {
	Execute(ctx context.Context, target job.Target, cmd ...string) error
}

type jobExecutor struct {
	runner.DockerRunner
	StartFunc func(context.Context)
	EndFunc   func(context.Context, error)
}

// Execute job
func (r *jobExecutor) Execute(ctx context.Context, target job.Target, cmd ...string) error {
	r.StartFunc(ctx)

	workDir, cleanup, err := target.Prepare()
	if err != nil {
		r.EndFunc(ctx, err)
		return errors.WithStack(err)
	}
	defer cleanup()

	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		errs <- r.DockerRunner.Run(timeout, workDir, docker.Tag(random.String(16, random.Lowercase)), cmd)
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		r.EndFunc(ctx, timeout.Err())
		return timeout.Err()
	case err := <-errs:
		r.EndFunc(ctx, err)
		return err
	}
}
