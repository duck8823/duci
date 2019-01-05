package docker_test

import (
	"fmt"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewBuildLog(t *testing.T) {
	// when
	got := docker.NewBuildLog(strings.NewReader("hello world"))

	// then
	if _, ok := got.(job.Log); !ok {
		t.Errorf("type assertion error.")
	}
}

func TestBuildLogger_ReadLine(t *testing.T) {
	// given
	now := time.Now()
	defer docker.SetNowFunc(func() time.Time {
		return now
	})()

	// and
	want := &job.LogLine{
		Timestamp: now,
		Message:   "hello test",
	}

	// and
	sut := docker.NewBuildLog(strings.NewReader(fmt.Sprintf("{\"stream\":\"%s\"}\n\nhello world\n", want.Message)))

	// when
	got, err := sut.ReadLine()

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but: %+v", cmp.Diff(got, want))
	}

	// when
	got, err = sut.ReadLine()

	// then
	if err != io.EOF {
		t.Errorf("error must be io.EOF, but got %+v", err)
	}

	// and
	if got != nil {
		t.Errorf("must be nil, but got %+v", got.Message)
	}
}

func TestNewRunLog(t *testing.T) {
	// when
	got := docker.NewRunLog(strings.NewReader("hello world"))

	// then
	if _, ok := got.(job.Log); !ok {
		t.Errorf("type assertion error.")
	}
}

func TestRunLogger_ReadLine(t *testing.T) {
	// given
	now := time.Now()
	defer docker.SetNowFunc(func() time.Time {
		return now
	})()

	// and
	want := &job.LogLine{
		Timestamp: now,
		Message:   "hello test",
	}

	// and
	sut := docker.NewRunLog(strings.NewReader(fmt.Sprintf("%shello test\rskipped line\n\n1234567890", string([]byte{1, 0, 0, 0, 9, 9, 9, 9}))))

	// when
	got, err := sut.ReadLine()

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but: %+v", cmp.Diff(got, want))
	}

	// when
	got, err = sut.ReadLine()

	// then
	if err == nil || err == io.EOF {
		t.Errorf("error must not be nil (invalid prefix), but got %+v", err)
	}

	// and
	if got != nil {
		t.Errorf("must be nil, but got %+v", got.Message)
	}

	// when
	got, err = sut.ReadLine()

	// then
	if err != io.EOF {
		t.Errorf("error must be io.EOF, but got %+v", err)
	}

	// and
	if got != nil {
		t.Errorf("must be nil, but got %+v", got.Message)
	}
}
