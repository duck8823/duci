package clock

import "time"

func SetDefaultFunc(f func() time.Time) {
	defaultFunc = f
}
