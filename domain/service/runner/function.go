package runner

import (
	"context"
	. "github.com/duck8823/duci/domain/model/job"
)

// LogFunc is function of Log
type LogFunc func(context.Context, Log)

// LogFuncs is slice of LogFunc
type LogFuncs []LogFunc

// Exec execute in goroutine
func (l LogFuncs) Exec(ctx context.Context, log Log) {
	for _, f := range l {
		go f(ctx, log)
	}
}
