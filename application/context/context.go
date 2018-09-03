package context

import (
	"context"
	"github.com/google/uuid"
	"net/url"
	"time"
)

type Context interface {
	context.Context
	UUID() uuid.UUID
	TaskName() string
	URL() *url.URL
}

type jobContext struct {
	context.Context
	uuid     uuid.UUID
	taskName string
	url      *url.URL
}

func New(taskName string, id uuid.UUID, url *url.URL) Context {
	return &jobContext{
		Context:  context.Background(),
		uuid:     id,
		taskName: taskName,
		url:      url,
	}
}

func (c *jobContext) UUID() uuid.UUID {
	return c.uuid
}

func (c *jobContext) TaskName() string {
	return c.taskName
}

func (c *jobContext) URL() *url.URL {
	return c.url
}

func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &jobContext{
		Context:  ctx,
		uuid:     parent.UUID(),
		taskName: parent.TaskName(),
		url:      parent.URL(),
	}, cancel
}
