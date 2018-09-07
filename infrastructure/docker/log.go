package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"time"
)

var now = time.Now

type Log interface {
	ReadLine() (*LogLine, error)
}

type LogLine struct {
	Timestamp time.Time
	Message   []byte
}

type buildLogger struct {
	reader *bufio.Reader
}

func (l *buildLogger) ReadLine() (*LogLine, error) {
	for {
		line, _, readErr := l.reader.ReadLine()
		msg := extractMessage(line)
		if readErr == io.EOF {
			return &LogLine{Timestamp: now(), Message: msg}, readErr
		}
		if readErr != nil {
			return nil, errors.WithStack(readErr)
		}

		if len(msg) == 0 {
			continue
		}

		return &LogLine{Timestamp: now(), Message: msg}, readErr
	}
}

type runLogger struct {
	reader *bufio.Reader
}

func (l *runLogger) ReadLine() (*LogLine, error) {
	for {
		line, _, readErr := l.reader.ReadLine()
		if readErr != nil && readErr != io.EOF {
			return nil, errors.WithStack(readErr)
		}

		messages, err := trimPrefix(line)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		// prevent to CR
		progress := bytes.Split(messages, []byte{'\r'})
		return &LogLine{Timestamp: now(), Message: progress[0]}, readErr
	}
}

func extractMessage(line []byte) []byte {
	s := &struct {
		Stream string `json:"stream"`
	}{}
	json.NewDecoder(bytes.NewReader(line)).Decode(s)
	return []byte(s.Stream)
}

func trimPrefix(line []byte) ([]byte, error) {
	if len(line) < 8 {
		return []byte{}, nil
	}

	// detect logstore prefix
	// see https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs
	if !((line[0] == 1 || line[0] == 2) && (line[1] == 0 && line[2] == 0 && line[3] == 0)) {
		return nil, fmt.Errorf("invalid logstore prefix: %+v", line[:7])
	}
	return line[8:], nil
}
