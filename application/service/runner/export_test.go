package runner

import (
	"github.com/duck8823/duci/infrastructure/docker"
	"io"
	"time"
)

type MockRepo struct {
	FullName string
	SSHURL   string
	CloneURL string
}

func (r *MockRepo) GetFullName() string {
	return r.FullName
}

func (r *MockRepo) GetSSHURL() string {
	return r.SSHURL
}

func (r *MockRepo) GetCloneURL() string {
	return r.CloneURL
}

type MockBuildLog struct {
}

func (l *MockBuildLog) ReadLine() (*docker.LogLine, error) {
	return &docker.LogLine{Timestamp: time.Now(), Message: []byte("{\"stream\":\"Hello World,\"}")}, io.EOF
}

type MockJobLog struct {
}

func (l *MockJobLog) ReadLine() (*docker.LogLine, error) {
	return &docker.LogLine{Timestamp: time.Now(), Message: []byte("Hello World,")}, io.EOF
}
