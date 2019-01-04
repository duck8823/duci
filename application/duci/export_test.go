package duci

import (
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"io"
	"net/url"
	"time"
)

type Duci = duci

func (d *Duci) SetJobService(service job_service.Service) (reset func()) {
	tmp := d.jobService
	d.jobService = service
	return func() {
		d.jobService = tmp
	}
}

func (d *Duci) SetGitHub(hub github.GitHub) (reset func()) {
	tmp := d.github
	d.github = hub
	return func() {
		d.github = tmp
	}
}

func URLMust(url *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	return url
}

type MockLog struct {
	Msgs []string
}

func (l *MockLog) ReadLine() (*job.LogLine, error) {
	if len(l.Msgs) == 0 {
		return nil, io.EOF
	}
	msg := l.Msgs[0]
	l.Msgs = l.Msgs[1:]
	return &job.LogLine{Timestamp: time.Now(), Message: msg}, nil
}

func String(val string) *string {
	return &val
}
