package runner

import "github.com/duck8823/duci/domain/model/docker"

type Builder = builder

func (b *builder) SetDocker(docker docker.Docker) (reset func()) {
	tmp := b.docker
	b.docker = docker
	return func() {
		b.docker = tmp
	}
}

func (b *builder) GetLogFunc() LogFunc {
	return b.logFunc
}

func (b *builder) SetLogFunc(logFunc LogFunc) (reset func()) {
	tmp := b.logFunc
	b.logFunc = logFunc
	return func() {
		b.logFunc = tmp
	}
}

type DockerRunnerImpl = dockerRunnerImpl

func (r *DockerRunnerImpl) SetDocker(docker docker.Docker) (reset func()) {
	tmp := r.docker
	r.docker = docker
	return func() {
		r.docker = tmp
	}
}

func (r *DockerRunnerImpl) GetLogFunc() LogFunc {
	return r.logFunc
}

func (r *DockerRunnerImpl) SetLogFunc(logFunc LogFunc) (reset func()) {
	tmp := r.logFunc
	r.logFunc = logFunc
	return func() {
		r.logFunc = tmp
	}
}

var CreateTarball = createTarball

var DockerfilePath = dockerfilePath

var ExportedRuntimeOptions = runtimeOptions

type StubLog struct {
}
