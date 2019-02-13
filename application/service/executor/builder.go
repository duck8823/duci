package executor

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/pkg/errors"
)

// Builder is an executor builder
type Builder struct {
	docker    docker.Docker
	logFunc   runner.LogFunc
	initFunc  func(context.Context)
	startFunc func(context.Context)
	endFunc   func(context.Context, error)
}

// DefaultExecutorBuilder create new Builder of docker runner
func DefaultExecutorBuilder() (*Builder, error) {
	docker, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Builder{
		docker:    docker,
		logFunc:   runner.NothingToDo,
		initFunc:  nothingToDoStart,
		startFunc: nothingToDoStart,
		endFunc:   nothingToDoEnd,
	}, nil
}

// LogFunc set a LogFunc
func (b *Builder) LogFunc(f runner.LogFunc) *Builder {
	b.logFunc = f
	return b
}

// InitFunc set a initFunc
func (b *Builder) InitFunc(f func(context.Context)) *Builder {
	b.initFunc = f
	return b
}

// StartFunc set a startFunc
func (b *Builder) StartFunc(f func(context.Context)) *Builder {
	b.startFunc = f
	return b
}

// EndFunc set a endFunc
func (b *Builder) EndFunc(f func(context.Context, error)) *Builder {
	b.endFunc = f
	return b
}

// Build returns a executor
func (b *Builder) Build() Executor {
	r := runner.DefaultDockerRunnerBuilder().
		LogFunc(b.logFunc).
		Build()

	return &jobExecutor{
		DockerRunner: r,
		InitFunc:     b.initFunc,
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
