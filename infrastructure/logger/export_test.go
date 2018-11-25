package logger

import "time"

func SetNowFunc(f func() time.Time) (reset func()) {
	tmp := now
	now = f
	return func() {
		now = tmp
	}
}
