package docker_test

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/labstack/gommon/random"
	"github.com/moby/moby/client"
	"io/ioutil"
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
	cli, err := docker.New()
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	t.Run("with correct archive", func(t *testing.T) {
		t.Run("without sub directory", func(t *testing.T) {
			// given
			tag := strings.ToLower(random.String(64))

			tar, err := os.Open("testdata/correct_archive.tar")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}

			// when
			if err := cli.Build(context.New("test/task"), tar, tag, "./Dockerfile"); err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			images := dockerImages(t)

			// then
			if !contains(images, fmt.Sprintf("%s:latest", tag)) {
				t.Errorf("docker images must contains. images: %+v, tag: %+v", images, tag)
			}
		})

		t.Run("with sub directory", func(t *testing.T) {
			// given
			tag := strings.ToLower(random.String(64))

			tar, err := os.Open("testdata/correct_archive_subdir.tar")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}

			// when
			if err := cli.Build(context.New("test/task"), tar, tag, ".duci/Dockerfile"); err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			images := dockerImages(t)

			// then
			if !contains(images, fmt.Sprintf("%s:latest", tag)) {
				t.Errorf("docker images must contains. images: %+v, tag: %+v", images, tag)
			}
		})
	})

	t.Run("with invalid archive", func(t *testing.T) {
		// given
		tag := strings.ToLower(random.String(64))

		tar, err := os.Open("testdata/invalid_archive.tar")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// expect
		if err := cli.Build(context.New("test/task"), tar, tag, "./Dockerfile"); err == nil {
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
			// given
			imagePull(t, "hello-world:latest")

			// when
			containerId, err := cli.Run(context.New("test/task"), opts, "hello-world")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			logs := containerLogsString(t, containerId)

			// then
			if !strings.Contains(logs, "Hello from Docker!") {
				t.Error("logs must contains `Hello from Docker!`")
			}
		})

		t.Run("with command", func(t *testing.T) {
			// given
			imagePull(t, "centos:latest")

			// when
			containerId, err := cli.Run(context.New("test/task"), opts, "centos", "echo", "Hello-world")
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}
			logs := containerLogsString(t, containerId)

			// then
			if strings.Contains(logs, "hello-world") {
				t.Errorf("logs must be equal `hello-world`. actual: %+v", logs)
			}
		})

		t.Run("with missing command", func(t *testing.T) {
			// given
			imagePull(t, "centos:latest")

			// expect
			if _, err := cli.Run(context.New("test/task"), opts, "centos", "missing_command"); err == nil {
				t.Error("error must occur")
			}
		})

		t.Run("when exit code is not zero", func(t *testing.T) {
			// given
			imagePull(t, "centos:latest")

			// expect
			if _, err := cli.Run(context.New("test/task"), opts, "centos", "false"); err != docker.Failure {
				t.Errorf("error must be docker.Failure, but got %+v", err)
			}
		})
	})

	t.Run("with environments", func(t *testing.T) {
		// given
		imagePull(t, "centos:latest")

		// and
		opts := docker.RuntimeOptions{
			Environments: docker.Environments{"ENV": "hello-world"},
		}

		// when
		containerId, err := cli.Run(context.New("test/task"), opts, "centos", "sh", "-c", "echo $ENV")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		logs := containerLogsString(t, containerId)

		// then
		if !strings.Contains(logs, "hello-world") {
			t.Errorf("logs must be equal `hello-world`. actual: %+v", logs)
		}
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
	if err := cli.Rm(context.New("test/task"), containerId); err != nil {
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
	if err := cli.Rmi(context.New("test/task"), tag); err != nil {
		t.Fatalf("error occured: %+v", err)
	}

	images := dockerImages(t)

	//then
	if contains(images, tag) {
		t.Errorf("images must not contains tag. images: %+v, tag: %+v", images, tag)
	}
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

func TestVolumes_ToMap(t *testing.T) {
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
			expected: map[string]struct{} {
				"/hoge/fuga:/hoge/hoge": {},
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

	images, err := cli.ImageList(context.New("test/task"), types.ImageListOptions{})
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

	containers, err := cli.ContainerList(context.New("test/task"), types.ContainerListOptions{})
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

	reader, err := cli.ContainerLogs(context.New("test/task"), containerId, types.ContainerLogsOptions{
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

	stream, err := cli.ImagePull(context.New("test/task"), ref, types.ImagePullOptions{})
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
	con, err := cli.ContainerCreate(context.New("test/task"), config, nil, nil, "")
	if err != nil {
		t.Fatalf("error occured. %+v", err)
		return ""
	}
	return con.ID
}
