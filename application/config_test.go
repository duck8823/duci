package application_test

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
	"time"
)

func TestConfiguration_String(t *testing.T) {
	// given
	conf := &application.Configuration{
		Server: &application.Server{
			WorkDir:      "/path/to/work_dir",
			SSHKeyPath:   "/path/to/ssh_key_path",
			DatabasePath: "path/to/databasePath",
			Port:         1234,
		},
		GitHub: &application.GitHub{
			APIToken: "github_api_token",
		},
		Job: &application.Job{
			Timeout:     60,
			Concurrency: 8,
		},
	}

	// and
	expected := fmt.Sprintf(
		"{\"server\":{\"workdir\":\"%s\",\"port\":%d,\"sshKeyPath\":\"%s\",\"databasePath\":\"%s\"},"+
			"\"github\":{\"apiToken\":\"***\"},\"job\":{\"timeout\":%d,\"concurrency\":%d}}",
		conf.Server.WorkDir,
		conf.Server.Port,
		conf.Server.SSHKeyPath,
		conf.Server.DatabasePath,
		conf.Job.Timeout,
		conf.Job.Concurrency,
	)

	// when
	actual := conf.String()

	// then
	if actual != expected {
		t.Errorf("find differences:\n %+v", cmp.Diff(actual, expected))
	}
}

func TestConfiguration_Set(t *testing.T) {
	t.Run("with correct config path", func(t *testing.T) {
		// given
		expected := &application.Configuration{
			Server: &application.Server{
				WorkDir:      "/path/to/workdir",
				Port:         8823,
				SSHKeyPath:   "/path/to/ssh_key",
				DatabasePath: "/path/to/database",
			},
			GitHub: &application.GitHub{
				APIToken: "github_api_token",
			},
			Job: &application.Job{
				Timeout:     300,
				Concurrency: 5,
			},
		}

		// when
		err := application.Config.Set("testdata/config.yml")

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		if !reflect.DeepEqual(application.Config, expected) {
			t.Errorf("find differences:\n %+v", cmp.Diff(application.Config, expected))
		}
	})

	t.Run("with wrong path", func(t *testing.T) {
		// expect
		if err := application.Config.Set("path/to/nothing"); err == nil {
			t.Error("error must occur")
		}
	})
}

func TestConfiguration_Addr(t *testing.T) {
	// given
	application.Config.Server.Port = 8823

	// when
	actual := application.Config.Addr()

	// then
	if actual != ":8823" {
		t.Errorf("addr should equal :8823, but got %s", actual)
	}
}

func TestConfiguration_Timeout(t *testing.T) {
	// given
	application.Config.Job.Timeout = 8823

	// when
	actual := application.Config.Timeout()

	// then
	if actual != 8823*time.Second {
		t.Errorf("addr should equal 8823 sec, but got %+v", actual)
	}
}
