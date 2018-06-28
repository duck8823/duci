package context

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Context interface {
	context.Context
	UUID() uuid.UUID
}

type contextWithUUID struct {
	context.Context
	uuid uuid.UUID
}

func New() Context {
	return &contextWithUUID{Context: context.Background(), uuid: uuid.New()}
}

func (c *contextWithUUID) UUID() uuid.UUID {
	return c.uuid
}

func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &contextWithUUID{ctx, parent.UUID()}, cancel
}
