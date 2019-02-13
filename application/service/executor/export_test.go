package executor

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
)

func (b *Builder) SetDocker(docker docker.Docker) (reset func()) {
	tmp := b.docker
	b.docker = docker
	return func() {
		b.docker = tmp
	}
}

func (b *Builder) SetInitFunc(initFunc func(context.Context)) (reset func()) {
	tmp := b.initFunc
	b.initFunc = initFunc
	return func() {
		b.initFunc = tmp
	}
}

func (b *Builder) SetStartFunc(startFunc func(context.Context)) (reset func()) {
	tmp := b.startFunc
	b.startFunc = startFunc
	return func() {
		b.startFunc = tmp
	}
}

func (b *Builder) SetLogFunc(logFunc func(context.Context, job.Log)) (reset func()) {
	tmp := b.logFunc
	b.logFunc = logFunc
	return func() {
		b.logFunc = tmp
	}
}

func (b *Builder) SetEndFunc(endFunc func(context.Context, error)) (reset func()) {
	tmp := b.endFunc
	b.endFunc = endFunc
	return func() {
		b.endFunc = tmp
	}
}

var NothingToDoStart = nothingToDoStart
var NothingToDoEnd = nothingToDoEnd

type JobExecutor = jobExecutor

func (r *JobExecutor) SetInitFunc(initFunc func(context.Context)) (reset func()) {
	tmp := r.InitFunc
	r.InitFunc = initFunc
	return func() {
		r.InitFunc = tmp
	}
}

func (r *JobExecutor) SetStartFunc(startFunc func(context.Context)) (reset func()) {
	tmp := r.StartFunc
	r.StartFunc = startFunc
	return func() {
		r.StartFunc = tmp
	}
}

func (r *JobExecutor) SetEndFunc(endFunc func(context.Context, error)) (reset func()) {
	tmp := r.EndFunc
	r.EndFunc = endFunc
	return func() {
		r.EndFunc = tmp
	}
}

func (r *JobExecutor) SetDockerRunner(runner runner.DockerRunner) (reset func()) {
	tmp := r.DockerRunner
	r.DockerRunner = runner
	return func() {
		r.DockerRunner = tmp
	}
}

type StubTarget struct {
	Dir     job.WorkDir
	Cleanup job.Cleanup
	Err     error
}

func (t *StubTarget) Prepare() (dir job.WorkDir, cleanup job.Cleanup, err error) {
	return t.Dir, t.Cleanup, t.Err
}
