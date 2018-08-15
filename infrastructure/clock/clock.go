package clock

import "time"

var Now = defaultFunc

var defaultFunc = func() time.Time {
	return time.Now()
}

func Adjust() {
	Now = defaultFunc
}
