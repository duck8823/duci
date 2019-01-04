package executor

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/pkg/errors"
)

type builder struct {
	docker    docker.Docker
	logFunc   runner.LogFunc
	startFunc func(context.Context)
	endFunc   func(context.Context, error)
}

// DefaultExecutorBuilder create new builder of docker runner
func DefaultExecutorBuilder() (*builder, error) {
	docker, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &builder{
		docker:    docker,
		logFunc:   runner.NothingToDo,
		startFunc: nothingToDoStart,
		endFunc:   nothingToDoEnd,
	}, nil
}

// LogFunc set a LogFunc
func (b *builder) LogFunc(f runner.LogFunc) *builder {
	b.logFunc = f
	return b
}

// StartFunc set a startFunc
func (b *builder) StartFunc(f func(context.Context)) *builder {
	b.startFunc = f
	return b
}

// EndFunc set a endFunc
func (b *builder) EndFunc(f func(context.Context, error)) *builder {
	b.endFunc = f
	return b
}

// Build returns a executor
func (b *builder) Build() *jobExecutor {
	r := runner.DefaultDockerRunnerBuilder().
		LogFunc(b.logFunc).
		Build()

	return &jobExecutor{
		DockerRunner: r,
		StartFunc:    b.startFunc,
		EndFunc:      b.endFunc,
	}
}

var nothingToDoStart = func(context.Context) {
	// nothing to do
}

var nothingToDoEnd = func(context.Context, error) {
	// nothing to do
}
