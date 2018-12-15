package github

import "net/url"

// State represents state of commit status
type State string

func (s State) String() string {
	return string(s)
}

type Description string

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

type CommitStatus struct {
	TargetSource *TargetSource
	State        State
	Description  Description
	Context      string
	TargetURL    *url.URL
}
