package model

type Job struct {
	Finished bool      `json:"finished"`
	Stream   []Message `json:"stream"`
}

type Message struct {
	Level string `json:"level"`
	Time  string `json:"time"`
	Text  string `json:"message"`
}
