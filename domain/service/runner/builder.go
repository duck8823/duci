package runner

import (
	"github.com/duck8823/duci/domain/service/docker"
)

// DefaultDockerRunnerBuilder create new builder of docker runner
func DefaultDockerRunnerBuilder() *builder {
	cli, _ := docker.New()
	return &builder{docker: cli}
}

type builder struct {
	docker   docker.Docker
	logFuncs LogFuncs
}

// LogFunc append a LogFunc
func (b *builder) LogFunc(f LogFunc) *builder {
	b.logFuncs = append(b.logFuncs, f)
	return b
}

// Build returns a docker runner
func (b *builder) Build() *dockerRunnerImpl {
	return &dockerRunnerImpl{
		Docker:   b.docker,
		LogFuncs: b.logFuncs,
	}
}
