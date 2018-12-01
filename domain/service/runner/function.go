package runner

import (
	"context"
	. "github.com/duck8823/duci/domain/model/job"
)

// LogFunc is function of Log
type LogFunc func(context.Context, Log)
