package semaphore_test

import (
	"context"
	"github.com/duck8823/duci/application/semaphore"
	"testing"
	"time"
)

func TestMake(t *testing.T) {
	// expect
	if err := semaphore.Make(); err != nil {
		t.Error("error must not occure at first")
	}

	if err := semaphore.Make(); err == nil {
		t.Error("errpr must occure")
	}
}

func TestSemaphore(t *testing.T) {
	// given
	end := make(chan struct{}, 1)

	// and
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// when
	go func() {
		semaphore.Acquire()
		semaphore.Release()
		end <- struct{}{}
	}()

	// then
	select {
	case <-timeout.Done():
		t.Errorf("error occurred: %+v", timeout.Err())
	case <-end:
		// nothing to do
	}
}
