package github

// TargetPoint represents a target point for clone.
type TargetPoint interface {
	GetRef() string
	GetHead() string
}

// SimpleTargetPoint is a simple implementation for TargetPoint
type SimpleTargetPoint struct {
	Ref string
	SHA string
}

// GetRef returns a Ref
func (s *SimpleTargetPoint) GetRef() string {
	return s.Ref
}

// GetHead returns a SHA
func (s *SimpleTargetPoint) GetHead() string {
	return s.SHA
}
