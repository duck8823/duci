package runner

import (
	"context"
	"github.com/duck8823/duci/domain/model/job"
)

// LogFunc is function of Log
type LogFunc func(context.Context, job.Log)

// NothingToDo is function nothing to do
var NothingToDo = func(_ context.Context, _ job.Log) {}
