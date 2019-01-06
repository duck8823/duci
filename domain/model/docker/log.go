package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

var now = time.Now

type buildLogger struct {
	reader *bufio.Reader
}

// NewBuildLog return a instance of Log.
func NewBuildLog(r io.Reader) job.Log {
	return &buildLogger{bufio.NewReader(r)}
}

// ReadLine returns LogLine.
func (l *buildLogger) ReadLine() (*job.LogLine, error) {
	for {
		line, _, err := l.reader.ReadLine()
		if err != nil {
			return nil, err
		}

		msg := extractMessage(line)
		if len(msg) == 0 {
			continue
		}

		return &job.LogLine{Timestamp: now(), Message: msg}, nil
	}
}

type runLogger struct {
	reader *bufio.Reader
}

// NewRunLog returns a instance of Log
func NewRunLog(r io.Reader) job.Log {
	return &runLogger{bufio.NewReader(r)}
}

// ReadLine returns LogLine.
func (l *runLogger) ReadLine() (*job.LogLine, error) {
	for {
		line, _, err := l.reader.ReadLine()
		if err != nil {
			return nil, err
		}

		msg, err := trimPrefix(line)
		if err != nil {
			return nil, errors.WithStack(err)
		} else if len(msg) == 0 {
			continue
		}

		// prevent to CR
		progress := bytes.Split(msg, []byte{'\r'})
		return &job.LogLine{Timestamp: now(), Message: string(progress[0])}, nil
	}
}

func extractMessage(line []byte) string {
	s := &struct {
		Stream string `json:"stream"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(line)).Decode(s); err != nil {
		logrus.Error(err)
	}
	return s.Stream
}

func trimPrefix(line []byte) ([]byte, error) {
	if len(line) < 8 {
		return []byte{}, nil
	}

	// detect prefix
	// see https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs
	if !((line[0] == 1 || line[0] == 2) && (line[1] == 0 && line[2] == 0 && line[3] == 0)) {
		return nil, fmt.Errorf("invalid prefix: %+v", line[:7])
	}
	return line[8:], nil
}
