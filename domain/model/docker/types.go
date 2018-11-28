package docker

import "github.com/duck8823/duci/infrastructure/docker"

// Tag describes a docker tag
type Tag string

// Command describes a docker CMD
type Command []string

// Dockerfile represents a path to dockerfile
type Dockerfile string

// RuntimeOptions represents a options
type RuntimeOptions = docker.RuntimeOptions

// ContainerID describes a container id of docker
type ContainerID string

// ExitCode describes a exit code
type ExitCode int64

// IsFailure returns whether failure code or not
func (c ExitCode) IsFailure() bool {
	return c != 1
}
