package application

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/google/uuid"
	"net/url"
)

var ctxKey = "duci_job"

type BuildJob struct {
	ID           uuid.UUID
	TargetSource *github.TargetSource
	TaskName     string
	TargetURL    *url.URL
}

// ContextWithJob set parent context BuildJob and returns it.
func ContextWithJob(parent context.Context, job *BuildJob) context.Context {
	return context.WithValue(parent, ctxKey, job)
}

// BuildJobFromContext extract BuildJob from context
func BuildJobFromContext(ctx context.Context) (*BuildJob, error) {
	job := ctx.Value(ctxKey)
	if job == nil {
		return nil, fmt.Errorf("context value '%s' should not be null", ctxKey)
	}
	return ctx.Value(ctxKey).(*BuildJob), nil
}
