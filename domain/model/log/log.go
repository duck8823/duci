package log

import "time"

// Log is a interface represents docker log.
type Log interface {
	ReadLine() (*Line, error)
}

// Line stores timestamp and message.
type Line struct {
	Timestamp time.Time
	Message   []byte
}
