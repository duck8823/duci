package health

import "github.com/duck8823/duci/domain/model/docker"

type Handler = handler

func (c *Handler) SetDocker(docker docker.Docker) (reset func()) {
	tmp := c.docker
	c.docker = docker
	return func() {
		c.docker = tmp
	}
}
