package job

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Job represents a task
type Job struct {
	Finished bool      `json:"finished"`
	Stream   []LogLine `json:"stream"`
}

// AppendLog append log line to stream
func (j *Job) AppendLog(line LogLine) {
	j.Stream = append(j.Stream, line)
}

// Finish set true to Finished
func (j *Job) Finish() {
	j.Finished = true
}

// ToBytes returns marshal byte slice
func (j *Job) ToBytes() ([]byte, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// NewJob returns unmarshal Job instance
func NewJob(data []byte) (*Job, error) {
	job := &Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// ID is the identifier of job
type ID uuid.UUID

// ToSlice returns slice value
func (i ID) ToSlice() []byte {
	return []byte(uuid.UUID(i).String())
}
