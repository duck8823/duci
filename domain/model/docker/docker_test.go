package docker_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/docker/mock_docker"
	. "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/github"
	"github.com/labstack/gommon/random"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("without environment variable", func(t *testing.T) {
		// given
		DOCKER_CERT_PATH := os.Getenv("DOCKER_CERT_PATH")
		_ = os.Setenv("DOCKER_CERT_PATH", "")
		defer func() {
			_ = os.Setenv("DOCKER_CERT_PATH", DOCKER_CERT_PATH)
		}()

		DOCKER_HOST := os.Getenv("DOCKER_HOST")
		_ = os.Setenv("DOCKER_HOST", "")
		defer func() {
			_ = os.Setenv("DOCKER_HOST", DOCKER_HOST)
		}()

		DOCKER_API_VERSION := os.Getenv("DOCKER_API_VERSION")
		_ = os.Setenv("DOCKER_API_VERSION", "")
		defer func() {
			_ = os.Setenv("DOCKER_API_VERSION", DOCKER_API_VERSION)
		}()

		// when
		get, err := docker.New()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		_, ok := get.(docker.Docker)
		if !ok {
			t.Error("instance must be docker.Docker")
		}
	})

	t.Run("with wrong environment variable DOCKER_HOST", func(t *testing.T) {
		// given
		DOCKER_HOST := os.Getenv("DOCKER_HOST")
		_ = os.Setenv("DOCKER_HOST", "wrong host name")
		defer func() {
			_ = os.Setenv("DOCKER_HOST", DOCKER_HOST)
		}()

		// when
		got, err := docker.New()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		if got != nil {
			t.Errorf("instance must be nil, but got %+v", err)
		}
	})
}

func TestClient_Build(t *testing.T) {
	t.Run("with collect ImageBuildResponse", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		// and
		ctx := context.Background()
		buildContext := strings.NewReader("hello world")
		tag := "test_tag"
		dockerfile := "testdata/Dockerfile"

		// and
		want := "want value"

		// and
		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ImageBuild(Eq(ctx), Eq(buildContext), Eq(types.ImageBuildOptions{
				Tags:       []string{tag},
				BuildArgs:  map[string]*string{},
				Dockerfile: dockerfile,
				Remove:     true,
			})).
			Times(1).
			Return(types.ImageBuildResponse{
				Body: ioutil.NewReadCloser(strings.NewReader(fmt.Sprintf("{\"stream\":\"%s\"}", want)), nil),
			}, nil)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// when
		got, err := sut.Build(ctx, buildContext, docker.Tag(tag), docker.Dockerfile{Dir: ".", Path: dockerfile})

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		line, _ := got.ReadLine()
		if line.Message != want {
			t.Errorf("want: %s, but got: %s", want, line.Message)
		}
	})

	t.Run("with build error", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		// and
		ctx := context.Background()
		buildContext := strings.NewReader("hello world")
		tag := "test_tag"
		dockerfile := "testdata/Dockerfile"

		// and
		empty := types.ImageBuildResponse{}
		wantError := errors.New("test_error")

		// and
		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ImageBuild(Eq(ctx), Eq(buildContext), Eq(types.ImageBuildOptions{
				Tags:       []string{tag},
				BuildArgs:  map[string]*string{},
				Dockerfile: dockerfile,
				Remove:     true,
			})).
			Times(1).
			Return(empty, wantError)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// when
		got, err := sut.Build(ctx, buildContext, docker.Tag(tag), docker.Dockerfile{Dir: ".", Path: dockerfile})

		// then
		if err.Error() != wantError.Error() {
			t.Errorf("error want: %+v, but got: %+v", wantError, err)
		}

		// and
		if got != nil {
			t.Errorf("log moust be nil, but got %+v", err)
		}
	})

	t.Run("when failure read log", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		// and
		ctx := context.Background()
		buildContext := strings.NewReader("hello world")
		tag := "test_tag"
		dockerfile := "testdata/Dockerfile"

		// and
		errorRes := types.ImageBuildResponse{
			 Body: new(docker.ErrorResponse),
		}

		// and
		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ImageBuild(Eq(ctx), Eq(buildContext), Eq(types.ImageBuildOptions{
				Tags:       []string{tag},
				BuildArgs:  map[string]*string{},
				Dockerfile: dockerfile,
				Remove:     true,
			})).
			Times(1).
			Return(errorRes, nil)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// when
		got, err := sut.Build(ctx, buildContext, docker.Tag(tag), docker.Dockerfile{Dir: ".", Path: dockerfile})

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("log moust be nil, but got %+v", err)
		}
	})

	t.Run("with invalid dockerfile path", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		// and
		ctx := context.Background()
		buildContext := strings.NewReader("hello world")
		tag := "test_tag"
		dockerfile := "invalid/testdata/Dockerfile"

		// and
		sut := &docker.Client{}

		// when
		got, err := sut.Build(ctx, buildContext, docker.Tag(tag), docker.Dockerfile{Dir: ".", Path: dockerfile})

		// then
		if err == nil {
			t.Errorf("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("log moust be nil, but got %+v", err)
		}
	})
}

func TestClient_Run(t *testing.T) {
	t.Run("nominal scenario", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		// and
		ctx := context.Background()
		opts := docker.RuntimeOptions{}
		tag := docker.Tag("test_tag")
		cmd := make([]string, 0)

		// and
		wantID := docker.ContainerID(random.String(16, random.Alphanumeric))
		want := "hello test"

		// and
		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ContainerCreate(Eq(ctx), Any(), Any(), Nil(), Eq("")).
			Times(1).
			Return(container.ContainerCreateCreatedBody{
				ID: wantID.String(),
			}, nil)

		mockMoby.EXPECT().
			ContainerStart(Eq(ctx), Eq(wantID.String()), Eq(types.ContainerStartOptions{})).
			Times(1).
			Return(nil)

		mockMoby.EXPECT().
			ContainerLogs(Eq(ctx), Eq(wantID.String()), Eq(types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     true,
			})).
			Times(1).
			Return(
				ioutil.NewReadCloser(bytes.NewReader(append([]byte{1, 0, 0, 0, 1, 1, 1, 1}, []byte(want)...)),
					nil,
				),
				nil,
			)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// when
		gotID, got, err := sut.Run(ctx, opts, docker.Tag(tag), cmd)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if gotID != wantID {
			t.Errorf("id want: %s, but got: %s", wantID, gotID)
		}

		// and
		line, _ := got.ReadLine()
		if line.Message != want {
			t.Errorf("want: %s, but got: %s", want, line.Message)
		}
	})

	t.Run("non-nominal scenarios", func(t *testing.T) {
		// where
		for _, tt := range []struct {
			name    string
			f       func(*testing.T, docker.RunArgs) (moby docker.Moby, finish func())
			emptyID bool
		}{
			{
				name:    "when failed create container",
				f:       mockMobyFailedCreateContainer,
				emptyID: true,
			},
			{
				name:    "when failed start container",
				f:       mockMobyFailedContainerStart,
				emptyID: false,
			},
			{
				name:    "when failed container logs",
				f:       mockMobyFailedContainerLogs,
				emptyID: false,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// given
				ctx := context.Background()
				opts := docker.RuntimeOptions{}
				tag := docker.Tag("test_tag")
				cmd := make([]string, 0)

				// and
				mockMoby, finish := tt.f(t, docker.RunArgs{Ctx: ctx, Opts: opts, Tag: tag, Cmd: cmd})
				defer finish()

				// and
				sut := &docker.Client{}
				defer sut.SetMoby(mockMoby)()

				// when
				gotID, got, err := sut.Run(ctx, opts, tag, cmd)

				// then
				if err == nil {
					t.Error("error must not be nil")
				}

				// and
				if tt.emptyID && len(gotID.String()) != 0 {
					t.Errorf("id must be empty, but got: %s", gotID)
				} else if !tt.emptyID && len(gotID.String()) == 0 {
					t.Error("id must not be empty")
				}

				// and
				if got != nil {
					t.Errorf("log must be nil, but got: %+v", got)
				}
			})
		}
	})
}

func mockMobyFailedCreateContainer(
	t *testing.T,
	args docker.RunArgs,
) (moby docker.Moby, finish func()) {
	t.Helper()

	ctrl := NewController(t)

	mockMoby := mock_docker.NewMockMoby(ctrl)
	mockMoby.EXPECT().
		ContainerCreate(Eq(args.Ctx), Any(), Any(), Nil(), Eq("")).
		Times(1).
		Return(container.ContainerCreateCreatedBody{
			ID: random.String(16, random.Alphanumeric),
		}, errors.New("test error"))

	mockMoby.EXPECT().
		ContainerStart(Any(), Any(), Any()).
		Times(0)

	mockMoby.EXPECT().
		ContainerLogs(Any(), Any(), Any()).
		Times(0)

	return mockMoby, func() {
		ctrl.Finish()
	}
}

func mockMobyFailedContainerStart(
	t *testing.T,
	args docker.RunArgs,
) (moby docker.Moby, finish func()) {
	t.Helper()

	conID := random.String(16, random.Alphanumeric)

	ctrl := NewController(t)

	mockMoby := mock_docker.NewMockMoby(ctrl)
	mockMoby.EXPECT().
		ContainerCreate(Eq(args.Ctx), Any(), Any(), Nil(), Eq("")).
		Times(1).
		Return(container.ContainerCreateCreatedBody{
			ID: conID,
		}, nil)

	mockMoby.EXPECT().
		ContainerStart(Eq(args.Ctx), Eq(conID), Eq(types.ContainerStartOptions{})).
		Times(1).
		Return(errors.New("test error"))

	mockMoby.EXPECT().
		ContainerLogs(Any(), Any(), Any()).
		Times(0)

	return mockMoby, func() {
		ctrl.Finish()
	}
}

func mockMobyFailedContainerLogs(
	t *testing.T,
	args docker.RunArgs,
) (moby docker.Moby, finish func()) {
	t.Helper()

	conID := random.String(16, random.Alphanumeric)

	ctrl := NewController(t)

	mockMoby := mock_docker.NewMockMoby(ctrl)
	mockMoby.EXPECT().
		ContainerCreate(Eq(args.Ctx), Any(), Any(), Nil(), Eq("")).
		Times(1).
		Return(container.ContainerCreateCreatedBody{
			ID: conID,
		}, nil)

	mockMoby.EXPECT().
		ContainerStart(Eq(args.Ctx), Eq(conID), Eq(types.ContainerStartOptions{})).
		Times(1).
		Return(nil)

	mockMoby.EXPECT().
		ContainerLogs(Eq(args.Ctx), Eq(conID), Eq(types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})).
		Times(1).
		Return(nil, errors.New("test error"))

	return mockMoby, func() {
		ctrl.Finish()
	}
}

func TestClient_RemoveContainer(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		// given
		ctx := context.Background()
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ContainerRemove(Eq(ctx), Eq(conID.String()), Eq(types.ContainerRemoveOptions{})).
			Times(1).
			Return(nil)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// expect
		if err := sut.RemoveContainer(ctx, conID); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctx := context.Background()
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ContainerRemove(Eq(ctx), Eq(conID.String()), Eq(types.ContainerRemoveOptions{})).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// expect
		if err := sut.RemoveContainer(ctx, conID); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestClient_RemoveImage(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		// given
		ctx := context.Background()
		tag := docker.Tag(random.String(16, random.Alphanumeric))

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ImageRemove(Eq(ctx), Eq(tag.String()), Eq(types.ImageRemoveOptions{})).
			Times(1).
			Return(nil, nil)

		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// expect
		if err := sut.RemoveImage(ctx, tag); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("without error", func(t *testing.T) {
		// given
		ctx := context.Background()
		tag := docker.Tag(random.String(16, random.Alphanumeric))

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ImageRemove(Eq(ctx), Eq(tag.String()), Eq(types.ImageRemoveOptions{})).
			Times(1).
			Return(nil, errors.New("test error"))

		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// expect
		if err := sut.RemoveImage(ctx, tag); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestClient_ExitCode(t *testing.T) {
	t.Run("with exit code", func(t *testing.T) {
		// given
		ctx := context.Background()
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		// and
		want := docker.ExitCode(19)

		// and
		body := make(chan container.ContainerWaitOKBody, 1)
		e := make(chan error, 1)

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ContainerWait(Eq(ctx), Eq(conID.String()), Eq(container.WaitConditionNotRunning)).
			Times(1).
			Return(body, e)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// and
		body <- container.ContainerWaitOKBody{StatusCode: int64(want)}

		// when
		got, err := sut.ExitCode(ctx, conID)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal but: %+v", cmp.Diff(got, want))
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctx := context.Background()
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		// and
		want := docker.ExitCode(-1)

		// and
		body := make(chan container.ContainerWaitOKBody, 1)
		e := make(chan error, 1)

		// and
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			ContainerWait(Eq(ctx), Eq(conID.String()), Eq(container.WaitConditionNotRunning)).
			Times(1).
			Return(body, e)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// and
		e <- errors.New("test error")

		// when
		got, err := sut.ExitCode(ctx, conID)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal but: %+v", cmp.Diff(got, want))
		}
	})
}

func TestClient_Status(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			Info(Any()).
			Times(1).
			Return(types.Info{}, nil)

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// want
		if err := sut.Status(); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctrl := NewController(t)
		defer ctrl.Finish()

		mockMoby := mock_docker.NewMockMoby(ctrl)
		mockMoby.EXPECT().
			Info(Any()).
			Times(1).
			Return(types.Info{}, errors.New("test error"))

		// and
		sut := &docker.Client{}
		defer sut.SetMoby(mockMoby)()

		// want
		if err := sut.Status(); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestBuildArgs(t *testing.T) {
	// given
	dockerfile := docker.Dockerfile{Dir: ".", Path: "testdata/Dockerfile"}

	// and
	hostArg2 := os.Getenv("ARGUMENT_2")
	_ = os.Setenv("ARGUMENT_2", "host_arg2")
	defer func() {
		_ = os.Setenv("ARGUMENT_2", hostArg2)
	}()

	hostArg5 := os.Getenv("ARGUMENT_5")
	_ = os.Setenv("ARGUMENT_5", "host_arg5")
	defer func() {
		_ = os.Setenv("ARGUMENT_5", hostArg5)
	}()

	// and
	want := map[string]*string{
		"ARGUMENT_2": github.String("host_arg2"),
		"ARGUMENT_5": github.String("host_arg5"),
	}

	// when
	got, err := docker.BuildArgs(dockerfile)

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal but: %+v", cmp.Diff(got, want))
	}
}
