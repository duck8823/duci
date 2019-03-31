package docker

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type Client = dockerImpl

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

type ErrorResponse struct {
}

func (e *ErrorResponse) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func (e *ErrorResponse) Close() error {
	return nil
}
