package clock_test

import (
	"github.com/duck8823/duci/infrastructure/clock"
	"testing"
	"time"
)

func TestAdjust(t *testing.T) {
	// given
	clock.SetDefaultFunc(func() time.Time {
		return time.Unix(0, 0)
	})

	// and
	clock.Now = func() time.Time {
		return time.Unix(1, 1)
	}

	// when
	clock.Adjust()

	// then
	actual := clock.Now()
	if actual != time.Unix(0, 0) {
		t.Errorf("wont %#v, but got %#v", time.Unix(0, 0), actual)
	}
}
