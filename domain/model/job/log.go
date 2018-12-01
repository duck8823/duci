package job

import "time"

// Log is a interface represents docker log.
type Log interface {
	ReadLine() (*LogLine, error)
}

// LogLine stores timestamp and message.
type LogLine struct {
	Timestamp time.Time
	Message   string
}
