package github

// State represents state of commit status
type State = string

const (
	// PENDING represents pending state.
	PENDING State = "pending"
	// SUCCESS represents success state.
	SUCCESS State = "success"
	// ERROR represents error state.
	ERROR State = "error"
	// FAILURE represents failure state.
	FAILURE State = "failure"
)
