package semaphore

import (
	"github.com/duck8823/duci/application"
	"github.com/pkg/errors"
	"runtime"
)

var (
	sem         = make(chan struct{}, runtime.NumCPU()) // default concurrency
	initialized = false
)

// Make create semaphore with configuration
func Make() error {
	if initialized {
		return errors.New("semaphore already created.")
	}
	sem = make(chan struct{}, application.Config.Server.Concurrency)
	initialized = true
	return nil
}

func Acquire() {
	sem <- struct{}{}
}

func Release() {
	<-sem
}
