package logger

import "time"

func SetNowFunc(f func() time.Time) {
	now = f
}
