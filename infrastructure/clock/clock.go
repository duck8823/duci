package clock

import "time"

var Now = defaultFunc

var defaultFunc = time.Now

func Adjust() {
	Now = defaultFunc
}
