package model

import "time"

type Job struct {
	Finished bool      `json:"finished"`
	Stream   []Message `json:"stream"`
}

type Message struct {
	Time time.Time `json:"time"`
	Text string    `json:"message"`
}
