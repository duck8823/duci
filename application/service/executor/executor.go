package executor

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/service/runner"
	"github.com/labstack/gommon/random"
)

type JobExecutor struct {
	runner.DockerRunner
	StartFunc []func(context.Context)
	EndFunc   []func(context.Context, error)
}

// Execute job
func (r *JobExecutor) Execute(ctx context.Context, dir string, cmd ...string) error {
	for _, f := range r.StartFunc {
		go f(ctx)
	}

	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		errs <- r.DockerRunner.Run(timeout, dir, docker.Tag(random.String(16, random.Alphanumeric)), cmd)
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		r.executeEndFunc(ctx, timeout.Err())
		return timeout.Err()
	case err := <-errs:
		r.executeEndFunc(ctx, err)
		return err
	}
}

// executeEndFunc execute functions
func (r *JobExecutor) executeEndFunc(ctx context.Context, err error) {
	for _, f := range r.EndFunc {
		go f(ctx, err)
	}
}
