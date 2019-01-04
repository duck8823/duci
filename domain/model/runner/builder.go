package runner

import (
	"github.com/duck8823/duci/domain/model/docker"
)

// Builder represents a builder of docker runner
type Builder struct {
	docker  docker.Docker
	logFunc LogFunc
}

// DefaultDockerRunnerBuilder create new builder of docker runner
func DefaultDockerRunnerBuilder() *Builder {
	cli, _ := docker.New()
	return &Builder{docker: cli, logFunc: NothingToDo}
}

// LogFunc append a LogFunc
func (b *Builder) LogFunc(f LogFunc) *Builder {
	b.logFunc = f
	return b
}

// Build returns a docker runner
func (b *Builder) Build() DockerRunner {
	return &dockerRunnerImpl{
		docker:  b.docker,
		logFunc: b.logFunc,
	}
}
