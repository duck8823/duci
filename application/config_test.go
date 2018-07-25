package application_test

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"reflect"
	"testing"
	"time"
)

func TestConfiguration_String(t *testing.T) {
	// given
	conf := &application.Configuration{
		Server: &application.Server{
			WorkDir:     "/path/to/work_dir",
			SSHKeyPath:  "/path/to/ssh_key_path",
			Port:        1234,
			Timeout:     60,
			Concurrency: 8,
		},
	}

	// and
	expected := fmt.Sprintf(
		"{\"server\":{\"workdir\":\"%s\",\"port\":%d,\"sshKeyPath\":\"%s\",\"timeout\":%d,\"concurrency\":%d}}",
		conf.Server.WorkDir,
		conf.Server.Port,
		conf.Server.SSHKeyPath,
		conf.Server.Timeout,
		conf.Server.Concurrency,
	)

	// when
	actual := conf.String()

	// then
	if actual != expected {
		t.Errorf("wont %s, but got %s", expected, actual)
	}
}

func TestConfiguration_Set(t *testing.T) {
	t.Run("with correct config path", func(t *testing.T) {
		// given
		expected := &application.Configuration{
			Server: &application.Server{
				WorkDir:     "/path/to/workdir",
				Port:        8823,
				SSHKeyPath:  "/path/to/ssh_key",
				Timeout:     600,
				Concurrency: 4,
			},
		}

		// when
		err := application.Config.Set("testdata/config.yml")

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		if !reflect.DeepEqual(application.Config, expected) {
			t.Errorf("wont %+v, but got %+v", expected, application.Config)
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
	application.Config.Server.Timeout = 8823

	// when
	actual := application.Config.Timeout()

	// then
	if actual != 8823*time.Second {
		t.Errorf("addr should equal 8823 sec, but got %+v", actual)
	}
}
