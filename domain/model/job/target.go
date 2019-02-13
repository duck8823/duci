package job

import "context"

// Target represents build target
type Target interface {
	Prepare(context.Context) (WorkDir, Cleanup, error)
}

// WorkDir is a working directory for build job
type WorkDir string

// String returns string value
func (w WorkDir) String() string {
	return string(w)
}

// Cleanup function for workdir
type Cleanup func()
