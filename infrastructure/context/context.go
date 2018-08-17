package context

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Context interface {
	context.Context
	UUID() uuid.UUID
	TaskName() string
	Host() string
}

type jobContext struct {
	context.Context
	uuid     uuid.UUID
	taskName string
	host     string
}

func New(taskName string, id uuid.UUID, host string) Context {
	return &jobContext{
		Context:  context.Background(),
		uuid:     id,
		taskName: taskName,
		host:     host,
	}
}

func (c *jobContext) UUID() uuid.UUID {
	return c.uuid
}

func (c *jobContext) TaskName() string {
	return c.taskName
}

func (c *jobContext) Host() string {
	return c.host
}

func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &jobContext{
		Context:  ctx,
		uuid:     parent.UUID(),
		taskName: parent.TaskName(),
		host:     parent.Host(),
	}, cancel
}
