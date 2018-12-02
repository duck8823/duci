package github

// TargetPoint represents a target point for clone.
type TargetPoint interface {
	GetRef() string
	GetHead() string
}
