package github

import "net/url"

// State represents state of commit status
type State string

// String returns string value
func (s State) String() string {
	return string(s)
}

// Description explain a commit status
type Description string

// TrimmedString returns length-fixed description
func (d Description) TrimmedString() string {
	if len(d) > 50 {
		return string([]rune(d)[:47]) + "..."
	}
	return string(d)
}

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

// CommitStatus represents a commit status
type CommitStatus struct {
	TargetSource *TargetSource
	State        State
	Description  Description
	Context      string
	TargetURL    *url.URL
}
