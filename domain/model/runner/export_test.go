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

var CreateTarball = createTarball

var DockerfilePath = dockerfilePath

var ExportedRuntimeOptions = runtimeOptions
