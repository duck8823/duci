package model

import "time"

// Job represents one of execution task.
type Job struct {
	Finished bool      `json:"finished"`
	Stream   []Message `json:"stream"`
}

// Message is a log of job.
type Message struct {
	Time time.Time `json:"time"`
	Text string    `json:"message"`
}
