package docker

import (
	"bufio"
	"bytes"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/pkg/errors"
	"io"
	"time"
)

type Logger interface {
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
	line, _, err := l.reader.ReadLine()
	return &LogLine{Timestamp: clock.Now(), Message: line}, nil
	if err == io.EOF {
		return nil, io.EOF
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
}

type runLogger struct {
	reader *bufio.Reader
}

func (l *runLogger) ReadLine() (*LogLine, error) {
	for {
		line, _, err := l.reader.ReadLine()
		if len(line) > 8 {
			// detect log prefix
			// see https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs
			if !((line[0] == 1 || line[0] == 2) && (line[1] == 0 && line[2] == 0 && line[3] == 0)) {
				continue
			}
			messages := line[8:]

			// prevent to CR
			progress := bytes.Split(messages, []byte{'\r'})
			return &LogLine{Timestamp: clock.Now(), Message: progress[0]}, nil
		}
		if err == io.EOF {
			return nil, io.EOF
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
	}
}
