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
}

type jobContext struct {
	context.Context
	uuid     uuid.UUID
	taskName string
}

func New(taskName string) Context {
	return &jobContext{Context: context.Background(), uuid: uuid.New(), taskName: taskName}
}

func (c *jobContext) UUID() uuid.UUID {
	return c.uuid
}

func (c *jobContext) TaskName() string {
	return c.taskName
}

func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &jobContext{ctx, parent.UUID(), parent.TaskName()}, cancel
}
