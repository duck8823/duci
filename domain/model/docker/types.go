package docker

import (
	"github.com/pkg/errors"
	"io"
	"os"
)

// Tag describes a docker tag
type Tag string

// ToString return string value
func (t Tag) String() string {
	return string(t)
}

// Command describes a docker CMD
type Command []string

// Slice returns slice values
func (c Command) Slice() []string {
	return []string(c)
}

// Dockerfile represents a path to dockerfile
type Dockerfile string

// ToString returns string value
func (d Dockerfile) String() string {
	return string(d)
}

// Open dockerfile
func (d Dockerfile) Open() (io.Reader, error) {
	file, err := os.Open(d.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return file, nil
}

// ContainerID describes a container id of docker
type ContainerID string

// ToString returns string value
func (c ContainerID) String() string {
	return string(c)
}

// ExitCode describes a exit code
type ExitCode int64

// IsFailure returns whether failure code or not
func (c ExitCode) IsFailure() bool {
	return c != 0
}
