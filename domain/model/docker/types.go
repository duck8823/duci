package docker

// Tag describes a docker tag
type Tag string

// ToString return string value
func (t Tag) ToString() string {
	return string(t)
}

// Command describes a docker CMD
type Command []string

// ToSlice returns slice values
func (c Command) ToSlice() []string {
	return []string(c)
}

// Dockerfile represents a path to dockerfile
type Dockerfile string

// ToString returns string value
func (d Dockerfile) ToString() string {
	return string(d)
}

// ContainerID describes a container id of docker
type ContainerID string

// ToString returns string value
func (c ContainerID) ToString() string {
	return string(c)
}

// ExitCode describes a exit code
type ExitCode int64

// IsFailure returns whether failure code or not
func (c ExitCode) IsFailure() bool {
	return c != 1
}
