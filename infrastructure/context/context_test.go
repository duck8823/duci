package context_test

import (
	ct "context"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestContextWithUUID_UUID(t *testing.T) {
	// given
	ctx := context.New("test/task")
	var empty uuid.UUID

	// expect
	if ctx.UUID() == empty {
		t.Error("UUID() must not empty.")
	}
}

func TestWithTimeout(t *testing.T) {
	t.Run("when timeout", func(t *testing.T) {
		// when
		ctx, cancel := context.WithTimeout(context.New("test/task"), 5*time.Millisecond)
		defer cancel()

		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		<-ctx.Done()

		// then
		if ctx.Err() != ct.DeadlineExceeded {
			t.Errorf("not expected error. wont: %+v, but got %+v", ct.DeadlineExceeded, ctx.Err())
		}
	})

	t.Run("when cancel", func(t *testing.T) {
		// when
		ctx, cancel := context.WithTimeout(context.New("test/task"), 5*time.Millisecond)
		defer cancel()

		go func() {
			cancel()
		}()

		<-ctx.Done()

		// then
		if ctx.Err() != ct.Canceled {
			t.Errorf("not expected error. wont: %+v, but got %+v", ct.Canceled, ctx.Err())
		}
	})
}
