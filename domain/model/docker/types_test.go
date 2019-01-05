package docker_test

import (
	"fmt"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/gommon/random"
	"testing"
)

func TestTag_String(t *testing.T) {
	// given
	want := "hello"

	// and
	sut := docker.Tag(want)

	// when
	got := sut.String()

	// then
	if got != want {
		t.Errorf("must equal: want %s, got %s", want, got)
	}
}

func TestCommand_Slice(t *testing.T) {
	// given
	want := []string{"test", "./..."}

	// and
	sut := docker.Command(want)

	// when
	got := sut.Slice()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must equal: want %+v, got %+v", want, got)
	}
}

func TestDockerfile_String(t *testing.T) {
	// given
	want := "duck8823/duci"

	// and
	sut := docker.Dockerfile(want)

	// when
	got := sut.String()

	// then
	if got != want {
		t.Errorf("must equal: want %s, got %s", want, got)
	}
}

func TestContainerID_String(t *testing.T) {
	// given
	want := random.String(16, random.Alphanumeric)

	// and
	sut := docker.ContainerID(want)

	// when
	got := sut.String()

	// then
	if got != want {
		t.Errorf("must equal: want %s, got %s", want, got)
	}
}

func TestExitCode_IsFailure(t *testing.T) {
	// where
	for _, tt := range []struct {
		code int64
		want bool
	}{
		{
			code: 0,
			want: false,
		},
		{
			code: -1,
			want: true,
		},
		{
			code: 1,
			want: true,
		},
	} {
		t.Run(fmt.Sprintf("when code is %+v", tt.code), func(t *testing.T) {
			// given
			sut := docker.ExitCode(tt.code)

			// expect
			if sut.IsFailure() != tt.want {
				t.Errorf("must be %+v, but got %+v", tt.want, sut.IsFailure())
			}
		})

	}
}
