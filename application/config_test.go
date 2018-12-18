package application_test

import (
	"github.com/duck8823/duci/application"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConfiguration_String(t *testing.T) {
	// given
	sut := &application.Configuration{}

	// expect
	if sut.String() != "" {
		t.Errorf("must be empty string, but got \"%s\"", sut.String())
	}
}

func TestConfiguration_Set(t *testing.T) {
	t.Run("with correct config path", func(t *testing.T) {
		// given
		expected := &application.Configuration{
			Server: &application.Server{
				WorkDir:      "/path/to/workdir",
				Port:         8823,
				DatabasePath: "/path/to/database",
			},
			GitHub: &application.GitHub{
				SSHKeyPath: "/path/to/ssh_key",
				APIToken:   "github_api_token",
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

	t.Run("parse environment variable", func(t *testing.T) {
		// given
		expected := "hello world"

		// and
		os.Setenv("TEST_CONF_ENV", expected)

		//
		err := application.Config.Set("testdata/config_with_env.yml")

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		if application.Config.Server.WorkDir != expected {
			t.Errorf("wont %+v, but got %+v", expected, application.Config.Server.WorkDir)
		}
	})

	t.Run("with wrong path", func(t *testing.T) {
		// expect
		if err := application.Config.Set("path/to/nothing"); err == nil {
			t.Error("error must occur")
		}
	})
}

func TestConfiguration_Type(t *testing.T) {
	// when
	actual := application.Config.Type()

	// expect
	if actual != "string" {
		t.Errorf("type should equal string, but got %+v", actual)
	}
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

func TestMaskString_MarshalJSON(t *testing.T) {
	// given
	sut := application.MaskString("hoge")

	// when
	actual, err := sut.MarshalJSON()

	// then
	if err != nil {
		t.Errorf("error must not occur, but got %+v", err)
	}

	if string(actual) != "\"***\"" {
		t.Errorf("wont masked string, but got '%s'", actual)
	}
}

func TestMaskString_String(t *testing.T) {
	// given
	want := "hoge"

	// and
	sut := application.MaskString(want)

	// when
	got := sut.String()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}
