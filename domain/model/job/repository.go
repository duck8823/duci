package job

import "errors"

var NotFound = errors.New("job not found")

// Repository is Job Repository
type Repository interface {
	FindBy(ID) (*Job, error)
	Save(Job) error
}
