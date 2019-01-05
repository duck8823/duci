package job

import "errors"

// ErrNotFound represents a job not found error
var ErrNotFound = errors.New("job not found")

// Repository is Job Repository
type Repository interface {
	FindBy(ID) (*Job, error)
	Save(Job) error
}
