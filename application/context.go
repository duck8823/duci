package application

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"net/url"
)

var ctxKey = "duci_job"

// BuildJob represents once of job
type BuildJob struct {
	ID           job.ID
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
	val := ctx.Value(ctxKey)
	if val == nil {
		return nil, fmt.Errorf("context value '%s' should not be null", ctxKey)
	}
	buildJob, ok := val.(*BuildJob)
	if !ok {
		return nil, fmt.Errorf("invalid type in context '%s'", ctxKey)
	}
	return buildJob, nil
}
