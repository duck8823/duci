package git_test

import (
	"bufio"
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
	"time"
)

func TestCloneLogger_ReadLine(t *testing.T) {
	// given
	now := time.Now()
	defer git.SetNowFunc(func() time.Time {
		return now
	})()

	// and
	want := &job.LogLine{
		Timestamp: now,
		Message:   "Hello World",
	}

	// and
	sut := &git.CloneLogger{}
	defer sut.SetReader(bufio.NewReader(strings.NewReader(fmt.Sprintf("%s\n\n", want.Message))))()

	// when
	got, err := sut.ReadLine()

	// then
	if err != nil {
		t.Errorf("err must be nil, but got %+v", err)
	}

	// and
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}

	// when
	got, err = sut.ReadLine()

	// then
	if err != io.EOF {
		t.Errorf("must be equal io.EOF, but got %+v", err)
	}
}

func TestProgressLogger_Write(t *testing.T) {
	// given
	var got string
	sut := &git.ProgressLogger{
		LogFunc: func(_ context.Context, log job.Log) {
			line, _ := log.ReadLine()
			got = line.Message
		},
	}
	sut.SetContext(context.Background())

	// and
	want := "hello world"

	// when
	_, err := sut.Write([]byte(want))

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	if got != want {
		t.Errorf("must be equal, but not\n%+v", cmp.Diff(got, want))
	}
}
