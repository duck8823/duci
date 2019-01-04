package job

import "github.com/duck8823/duci/domain/model/job"

// Service represents job service
type Service interface {
	FindBy(id job.ID) (*job.Job, error)
	Start(id job.ID) error
	Append(id job.ID, line job.LogLine) error
	Finish(id job.ID) error
}
