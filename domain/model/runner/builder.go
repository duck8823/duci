package runner

import (
	"github.com/duck8823/duci/domain/model/docker"
)

type builder struct {
	docker  docker.Docker
	logFunc LogFunc
}

// DefaultDockerRunnerBuilder create new builder of docker runner
func DefaultDockerRunnerBuilder() *builder {
	cli, _ := docker.New()
	return &builder{docker: cli, logFunc: NothingToDo}
}

// LogFunc append a LogFunc
func (b *builder) LogFunc(f LogFunc) *builder {
	b.logFunc = f
	return b
}

// Build returns a docker runner
func (b *builder) Build() *dockerRunnerImpl {
	return &dockerRunnerImpl{
		docker:  b.docker,
		logFunc: b.logFunc,
	}
}
