package docker

import (
	"context"
	"time"
)

type Client = client

func (c *Client) SetMoby(moby Moby) (reset func()) {
	tmp := c.moby
	c.moby = moby
	return func() {
		c.moby = tmp
	}
}

type RunArgs struct {
	Ctx  context.Context
	Opts RuntimeOptions
	Tag  Tag
	Cmd  Command
}

func SetNowFunc(f func() time.Time) (reset func()) {
	tmp := now
	now = f
	return func() {
		now = tmp
	}
}
