package docker_test

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/docker/mock_docker"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with wrong docker environment", func(t *testing.T) {
		// given
		dockerHost := os.Getenv("DOCKER_HOST")
		os.Setenv("DOCKER_HOST", "hoge")

		// expect
		if _, err := docker.New(); err == nil {
			t.Errorf("error must occur")
		}

		// cleanup
		os.Setenv("DOCKER_HOST", dockerHost)
	})
}

func TestClientImpl_Build(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("when success image build", func(t *testing.T) {
		// given
		expected := "hello world"
		sr := strings.NewReader(fmt.Sprintf("{\"stream\":\"%s\"}", expected))
		r := ioutil.NopCloser(sr)

		// and
		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(types.ImageBuildResponse{Body: r}, nil)

		sut.SetMoby(mockMoby)

		// when
		log, err := sut.Build(context.New("test/task", uuid.New(), &url.URL{}), nil, "", "")

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// and
		line, err := log.ReadLine()
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// and
		if string(line.Message) != expected {
			t.Errorf("must be equal. wont %#v, but got %#v", expected, string(line.Message))
		}
	})

	t.Run("when failure image build", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(types.ImageBuildResponse{}, errors.New("test error"))

		sut.SetMoby(mockMoby)

		// expect
		if _, err := sut.Build(
			context.New("test/task", uuid.New(), &url.URL{}),
			nil,
			"",
			"",
		); err == nil {
			t.Errorf("error must occur, but got %+v", err)
		}
	})
}

func TestClientImpl_Run(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("when failure create container", func(t *testing.T) {
		// given
		id := random.String(64, random.Alphanumeric, random.Symbols)

		// and
		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(container.ContainerCreateCreatedBody{ID: id}, errors.New("test error"))

		sut.SetMoby(mockMoby)

		// when
		actual, _, err := sut.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "hello-world")

		// then
		if actual != "" {
			t.Errorf("id must be empty string, but got %+v", actual)
		}

		if err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("when success create container", func(t *testing.T) {
		t.Run("when failure start container", func(t *testing.T) {
			// given
			id := random.String(64, random.Alphanumeric, random.Symbols)

			// and
			ctrl := gomock.NewController(t)
			mockMoby := mock_docker.NewMockMoby(ctrl)

			mockMoby.EXPECT().
				ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(container.ContainerCreateCreatedBody{ID: id}, nil)

			mockMoby.EXPECT().
				ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(errors.New("test error"))

			sut.SetMoby(mockMoby)

			// when
			actual, _, err := sut.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "hello-world")

			// then
			if actual != id {
				t.Errorf("id must be equal %+v, but got %+v", id, actual)
			}

			if err == nil {
				t.Error("error must occur, but got nil")
			}
		})

		t.Run("when success start container", func(t *testing.T) {
			t.Run("when failure get log", func(t *testing.T) {
				// given
				id := random.String(64, random.Alphanumeric, random.Symbols)

				// and
				ctrl := gomock.NewController(t)
				mockMoby := mock_docker.NewMockMoby(ctrl)

				mockMoby.EXPECT().
					ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(container.ContainerCreateCreatedBody{ID: id}, nil)

				mockMoby.EXPECT().
					ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(nil)

				mockMoby.EXPECT().
					ContainerLogs(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("test error"))

				sut.SetMoby(mockMoby)

				// when
				actual, _, err := sut.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "hello-world")

				// then
				if actual != id {
					t.Errorf("id must be equal %+v, but got %+v", id, actual)
				}

				if err == nil {
					t.Error("error must occur, but got nil")
				}
			})

			t.Run("when success get log", func(t *testing.T) {
				t.Run("with valid log", func(t *testing.T) {
					// given
					id := random.String(64, random.Alphanumeric, random.Symbols)

					prefix := []byte{1, 0, 0, 0, 1, 1, 1, 1}
					msg := "hello test"
					log := ioutil.NopCloser(bytes.NewReader(append(prefix, []byte(msg)...)))

					// and
					ctrl := gomock.NewController(t)
					mockMoby := mock_docker.NewMockMoby(ctrl)

					mockMoby.EXPECT().
						ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						AnyTimes().
						Return(container.ContainerCreateCreatedBody{ID: id}, nil)

					mockMoby.EXPECT().
						ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).
						AnyTimes().
						Return(nil)

					mockMoby.EXPECT().
						ContainerLogs(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(log, nil)

					sut.SetMoby(mockMoby)

					// when
					actualID, actualLog, err := sut.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "hello-world")

					// then
					if actualID != id {
						t.Errorf("id must be equal %+v, but got %+v", id, actualID)
					}

					if actualLog == nil {
						t.Errorf("log must not nil")
					} else {
						line, _ := actualLog.ReadLine()
						if string(line.Message) != msg {
							t.Errorf("message must equal. wont %+v, but got %+v", msg, string(line.Message))
						}
					}

					if err != nil {
						t.Error("error must occur, but got nil")
					}
				})
			})
		})
	})
}

func TestClientImpl_Rm(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("when success removing container", func(t *testing.T) {
		// given
		conID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ContainerRemove(gomock.Any(), gomock.Eq(conID), gomock.Any()).
			Return(nil)

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rm(context.New("test/task", uuid.New(), &url.URL{}), conID); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("when failure removing container", func(t *testing.T) {
		// given
		conID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ContainerRemove(gomock.Any(), gomock.Eq(conID), gomock.Any()).
			Return(errors.New("test error"))

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rm(context.New("test/task", uuid.New(), &url.URL{}), conID); err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestClientImpl_Rmi(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("when success removing image", func(t *testing.T) {
		// given
		imageID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageRemove(gomock.Any(), gomock.Eq(imageID), gomock.Any()).
			Return(nil, nil)

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rmi(context.New("test/task", uuid.New(), &url.URL{}), imageID); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("when failure removing image", func(t *testing.T) {
		// given
		imageID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageRemove(gomock.Any(), gomock.Eq(imageID), gomock.Any()).
			Return(nil, errors.New("test error"))

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rmi(context.New("test/task", uuid.New(), &url.URL{}), imageID); err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestClientImpl_ExitCode2(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("when success removing image", func(t *testing.T) {
		// given
		imageID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageRemove(gomock.Any(), gomock.Eq(imageID), gomock.Any()).
			Return(nil, nil)

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rmi(context.New("test/task", uuid.New(), &url.URL{}), imageID); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("when failure removing image", func(t *testing.T) {
		// given
		imageID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ImageRemove(gomock.Any(), gomock.Eq(imageID), gomock.Any()).
			Return(nil, errors.New("test error"))

		sut.SetMoby(mockMoby)

		// expect
		if err := sut.Rmi(context.New("test/task", uuid.New(), &url.URL{}), imageID); err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestClientImpl_ExitCode(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("with exit code 0", func(t *testing.T) {
		// given
		exitCode := int64(19)

		body := make(chan container.ContainerWaitOKBody, 1)
		err := make(chan error, 1)

		// and
		body <- container.ContainerWaitOKBody{StatusCode: exitCode}

		// and
		conID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ContainerWait(gomock.Any(), gomock.Eq(conID), gomock.Any()).
			Return(body, err)

		sut.SetMoby(mockMoby)

		// when
		if code, _ := sut.ExitCode(context.New("test/task", uuid.New(), &url.URL{}), conID); code != exitCode {
			t.Errorf("code must equal %+v, but got %+v", exitCode, code)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		body := make(chan container.ContainerWaitOKBody, 1)
		err := make(chan error, 1)

		// and
		err <- errors.New("test error")

		// and
		conID := random.String(16, random.Alphanumeric, random.Symbols)

		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			ContainerWait(gomock.Any(), gomock.Eq(conID), gomock.Any()).
			Return(body, err)

		sut.SetMoby(mockMoby)

		// when
		if _, actualErr := sut.ExitCode(context.New("test/task", uuid.New(), &url.URL{}), conID); actualErr == nil {
			t.Error("error must occur but got nil")
		}
	})
}

func TestClientImpl_Info(t *testing.T) {
	// setup
	sut, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	t.Run("without error", func(t *testing.T) {
		// given
		expected := types.Info{ID: uuid.New().String()}

		// and
		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			Info(gomock.Any()).
			Return(expected, nil)

		sut.SetMoby(mockMoby)

		// when
		actual, err := sut.Info(context.New("test", uuid.New(), nil))

		// then
		if !cmp.Equal(actual, expected) {
			t.Errorf("must be equal. %+v", cmp.Diff(actual, expected))
		}

		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		mockMoby := mock_docker.NewMockMoby(ctrl)

		mockMoby.EXPECT().
			Info(gomock.Any()).
			Return(types.Info{}, errors.New("test"))

		sut.SetMoby(mockMoby)

		// expect
		if _, err := sut.Info(context.New("test", uuid.New(), nil)); err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestEnvironments_ToArray(t *testing.T) {
	var empty []string
	for _, testcase := range []struct {
		in       docker.Environments
		expected []string
	}{
		{
			in:       docker.Environments{},
			expected: empty,
		},
		{
			in: docker.Environments{
				"int":    19,
				"string": "hello",
			},
			expected: []string{
				"int=19",
				"string=hello",
			},
		},
	} {
		// when
		actual := testcase.in.ToArray()
		expected := testcase.expected
		sort.Strings(actual)
		sort.Strings(expected)

		// then
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("must be equal. actual=%+v, wont=%+v", actual, expected)
		}
	}
}

func TestVolumes_Volumes(t *testing.T) {
	for _, testcase := range []struct {
		in       docker.Volumes
		expected map[string]struct{}
	}{
		{
			in:       docker.Volumes{},
			expected: make(map[string]struct{}),
		},
		{
			in: docker.Volumes{
				"/hoge/fuga:/hoge/hoge",
			},
			expected: map[string]struct{}{
				"/hoge/fuga": {},
			},
		},
	} {
		// when
		actual := testcase.in.ToMap()
		expected := testcase.expected

		// then
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("must be equal. actual=%+v, wont=%+v", actual, expected)
		}
	}
}
