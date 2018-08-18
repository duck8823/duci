package docker_test

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
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
	cli, err := docker.New()
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	t.Run("with correct archive", func(t *testing.T) {
		t.Run("without sub directory", func(t *testing.T) {
			t.Parallel()

			// given
			tag := strings.ToLower(random.String(64))

			tar, err := os.Open("testdata/correct_archive.tar")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}

			// when
			logger, err := cli.Build(context.New("test/task", uuid.New(), &url.URL{}), tar, tag, "./Dockerfile")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			wait(t, logger)

			images := dockerImages(t)

			// then
			if !contains(images, fmt.Sprintf("%s:latest", tag)) {
				t.Errorf("docker images must contains. images: %+v, tag: %+v", images, tag)
			}
		})

		t.Run("with sub directory", func(t *testing.T) {
			t.Parallel()

			// given
			tag := strings.ToLower(random.String(64))

			tar, err := os.Open("testdata/correct_archive_subdir.tar")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}

			// when
			logger, err := cli.Build(context.New("test/task", uuid.New(), &url.URL{}), tar, tag, ".duci/Dockerfile")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			wait(t, logger)

			images := dockerImages(t)

			// then
			if !contains(images, fmt.Sprintf("%s:latest", tag)) {
				t.Errorf("docker images must contains. images: %+v, tag: %+v", images, tag)
			}
		})
	})

	t.Run("with invalid archive", func(t *testing.T) {
		t.Parallel()

		// given
		tag := strings.ToLower(random.String(64))

		tar, err := os.Open("testdata/invalid_archive.tar")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// expect
		if _, err := cli.Build(context.New("test/task", uuid.New(), &url.URL{}), tar, tag, "./Dockerfile"); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestClientImpl_Run(t *testing.T) {
	// setup
	cli, err := docker.New()
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	t.Run("without environments", func(t *testing.T) {
		// setup
		opts := docker.RuntimeOptions{}

		t.Run("without command", func(t *testing.T) {
			t.Parallel()

			// given
			imagePull(t, "hello-world:latest")

			// when
			containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), opts, "hello-world")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			containerWait(t, containerId)

			logs := containerLogsString(t, containerId)

			// then
			if !strings.Contains(logs, "Hello from Docker!") {
				t.Error("logs must contains `Hello from Docker!`")
			}

			// cleanup
			removeContainer(t, containerId)
		})

		t.Run("with command", func(t *testing.T) {
			t.Parallel()

			// given
			imagePull(t, "centos:latest")

			// when
			containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), opts, "centos", "echo", "Hello-world")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			containerWait(t, containerId)

			logs := containerLogsString(t, containerId)

			// then
			if strings.Contains(logs, "hello-world") {
				t.Errorf("logs must be equal `hello-world`. actual: %+v", logs)
			}

			// cleanup
			removeContainer(t, containerId)
		})

		t.Run("with missing command", func(t *testing.T) {
			t.Parallel()

			// given
			imagePull(t, "centos:latest")

			// expect
			if _, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), opts, "centos", "missing_command"); err == nil {
				t.Error("error must occur")
			}
		})
	})

	t.Run("with environments", func(t *testing.T) {
		t.Parallel()

		// given
		imagePull(t, "centos:latest")

		// and
		opts := docker.RuntimeOptions{
			Environments: docker.Environments{"ENV": "hello-world"},
		}

		// when
		containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), opts, "centos", "sh", "-c", "echo hello $ENV")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		containerWait(t, containerId)

		logs := containerLogsString(t, containerId)

		// then
		if !strings.Contains(logs, "hello-world") {
			t.Errorf("logs must be equal `hello-world`. actual: %+v", logs)
		}

		// cleanup
		removeContainer(t, containerId)
	})

	t.Run("with volumes", func(t *testing.T) {
		if os.Getenv("CI") == "duci" {
			t.Skip("skip if CI ( Docker in Docker )")
			// TODO reduce external dependencies
		}
		t.Parallel()

		// given
		imagePull(t, "centos:latest")

		// and
		path, err := filepath.Abs("testdata")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		opts := docker.RuntimeOptions{
			Volumes: docker.Volumes{fmt.Sprintf("%s:/tmp/testdata", path)},
		}

		// when
		containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), opts, "centos", "cat", "/tmp/testdata/data")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		containerWait(t, containerId)

		logs := containerLogsString(t, containerId)

		// then
		if !strings.Contains(logs, "hello-world") {
			t.Errorf("logs must be equal `hello-world`. actual: %+v", logs)
		}

		// cleanup
		removeContainer(t, containerId)
	})
}

func TestClientImpl_Rm(t *testing.T) {
	// setup
	cli, err := docker.New()
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	// given
	tag := "alpine:3.5"
	imagePull(t, tag)
	containerId := containerCreate(t, tag)

	// when
	if err := cli.Rm(context.New("test/task", uuid.New(), &url.URL{}), containerId); err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	containers := dockerContainers(t)

	// then
	if contains(containers, tag) {
		t.Errorf("containers must not contains id. containers: %+v, tag: %+v", containers, containerId)
	}
}

func TestClientImpl_Rmi(t *testing.T) {
	// setup
	cli, err := docker.New()
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	// given
	tag := "alpine:2.6"
	imagePull(t, tag)

	// when
	if err := cli.Rmi(context.New("test/task", uuid.New(), &url.URL{}), tag); err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	images := dockerImages(t)

	//then
	if contains(images, tag) {
		t.Errorf("images must not contains tag. images: %+v, tag: %+v", images, tag)
	}
}

func TestClientImpl_ExitCode(t *testing.T) {
	t.Run("with exit code 0", func(t *testing.T) {
		// given
		cli, err := docker.New()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		imagePull(t, "alpine:latest")

		// and
		containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "alpine", "sh", "-c", "exit 0")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// when
		code, err := cli.ExitCode(context.New("test/task", uuid.New(), &url.URL{}), containerId)

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		if code != 0 {
			t.Errorf("not equal wont 0, but got %d", code)
		}

		// cleanup
		removeContainer(t, containerId)
	})

	t.Run("with exit code 1", func(t *testing.T) {
		// given
		cli, err := docker.New()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		imagePull(t, "alpine:latest")

		// and
		containerId, _, err := cli.Run(context.New("test/task", uuid.New(), &url.URL{}), docker.RuntimeOptions{}, "alpine", "sh", "-c", "exit 1")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// when
		code, err := cli.ExitCode(context.New("test/task", uuid.New(), &url.URL{}), containerId)

		// then
		if err != nil {
			t.Error("error must occur, but got nil")
		}

		if code != 1 {
			t.Errorf("not equal wont 1, but got %d", code)
		}

		// cleanup
		removeContainer(t, containerId)
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

func contains(strings []string, str string) bool {
	for _, s := range strings {
		if s == str {
			return true
		}
	}
	return false
}

func dockerImages(t *testing.T) []string {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	images, err := cli.ImageList(context.New("test/task", uuid.New(), &url.URL{}), types.ImageListOptions{})
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	var names []string
	for _, image := range images {
		names = append(names, image.RepoTags...)
	}

	return names
}

func dockerContainers(t *testing.T) []string {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	containers, err := cli.ContainerList(context.New("test/task", uuid.New(), &url.URL{}), types.ContainerListOptions{})
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	var ids []string
	for _, con := range containers {
		ids = append(ids, con.ID)
	}
	return ids
}

func containerLogsString(t *testing.T, containerId string) string {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	reader, err := cli.ContainerLogs(context.New("test/task", uuid.New(), &url.URL{}), containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	log, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	return string(log)
}

func imagePull(t *testing.T, ref string) {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	stream, err := cli.ImagePull(context.New("test/task", uuid.New(), &url.URL{}), ref, types.ImagePullOptions{})
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	// wait until pull
	if _, err := ioutil.ReadAll(stream); err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	images := dockerImages(t)
	if !contains(images, ref) {
		t.Fatalf("docker images must be contains %s", ref)
	}
}

func containerCreate(t *testing.T, ref string) string {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	config := &container.Config{
		Image: ref,
		Cmd:   []string{"hello", "world"},
	}
	con, err := cli.ContainerCreate(context.New("test/task", uuid.New(), &url.URL{}), config, nil, nil, "")
	if err != nil {
		t.Fatalf("error occured. %+v", err)
		return ""
	}
	return con.ID
}

func containerWait(t *testing.T, containerId string) {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	body, err2 := cli.ContainerWait(context.New("test/task", uuid.New(), &url.URL{}), containerId, container.WaitConditionNotRunning)
	select {
	case <-body:
		return
	case <-err2:
		t.Fatalf("error occured. %+v", err)
	}
}

func removeContainer(t *testing.T, containerId string) {
	t.Helper()

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	if err := cli.ContainerRemove(context.New("test/task", uuid.New(), &url.URL{}), containerId, types.ContainerRemoveOptions{}); err != nil {
		t.Fatalf("error occured. %+v", err)
	}
}

func wait(t *testing.T, logger docker.Log) {
	t.Helper()

	for {
		_, err := logger.ReadLine()
		if err != nil {
			break
		}
	}
}
