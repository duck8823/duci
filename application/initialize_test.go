package application_test

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/labstack/gommon/random"
	"os"
	"path"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("when singleton container is empty", func(t *testing.T) {
		t.Run("without ssh key path", func(t *testing.T) {
			// given
			sshKeyPath := application.Config.GitHub.SSHKeyPath
			databasePath := application.Config.Server.DatabasePath
			application.Config.GitHub.SSHKeyPath = ""
			application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
			defer func() {
				application.Config.GitHub.SSHKeyPath = sshKeyPath
				application.Config.Server.DatabasePath = databasePath
			}()

			// and
			container.Clear()

			// when
			err := application.Initialize()

			// then
			if err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			git := new(git.Git)
			if err := container.Get(git); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			github := new(github.GitHub)
			if err := container.Get(github); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			jobService := new(job.Service)
			if err := container.Get(jobService); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}
		})

		t.Run("with correct ssh key path", func(t *testing.T) {
			// given
			dir := path.Join(os.TempDir(), random.String(16))
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			keyPath := path.Join(dir, "id_rsa")
			application.GenerateSSHKey(t, keyPath)

			// and
			sshKeyPath := application.Config.GitHub.SSHKeyPath
			databasePath := application.Config.Server.DatabasePath
			application.Config.GitHub.SSHKeyPath = keyPath
			application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
			defer func() {
				application.Config.GitHub.SSHKeyPath = sshKeyPath
				application.Config.Server.DatabasePath = databasePath
			}()

			// and
			container.Clear()

			// when
			err := application.Initialize()

			// then
			if err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			git := new(git.Git)
			if err := container.Get(git); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			github := new(github.GitHub)
			if err := container.Get(github); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			jobService := new(job.Service)
			if err := container.Get(jobService); err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}
		})

		t.Run("with invalid key path", func(t *testing.T) {
			// given
			sshKeyPath := application.Config.GitHub.SSHKeyPath
			databasePath := application.Config.Server.DatabasePath
			application.Config.GitHub.SSHKeyPath = "/path/to/invalid/key/path"
			application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
			defer func() {
				application.Config.GitHub.SSHKeyPath = sshKeyPath
				application.Config.Server.DatabasePath = databasePath
			}()

			// and
			container.Clear()

			// when
			err := application.Initialize()

			// then
			if err == nil {
				t.Error("error must not be nil")
			}
		})
	})

	t.Run("when singleton container contains Git instance", func(t *testing.T) {
		// given
		sshKeyPath := application.Config.GitHub.SSHKeyPath
		databasePath := application.Config.Server.DatabasePath
		application.Config.GitHub.SSHKeyPath = ""
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			application.Config.GitHub.SSHKeyPath = sshKeyPath
			application.Config.Server.DatabasePath = databasePath
		}()

		// and
		container.Override(new(git.Git))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when singleton container contains GitHub instance", func(t *testing.T) {
		// given
		sshKeyPath := application.Config.GitHub.SSHKeyPath
		databasePath := application.Config.Server.DatabasePath
		application.Config.GitHub.SSHKeyPath = ""
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			application.Config.GitHub.SSHKeyPath = sshKeyPath
			application.Config.Server.DatabasePath = databasePath
		}()

		// and
		container.Override(new(github.GitHub))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when singleton container contains JobService instance", func(t *testing.T) {
		// given
		sshKeyPath := application.Config.GitHub.SSHKeyPath
		databasePath := application.Config.Server.DatabasePath
		application.Config.GitHub.SSHKeyPath = ""
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			application.Config.GitHub.SSHKeyPath = sshKeyPath
			application.Config.Server.DatabasePath = databasePath
		}()

		// and
		container.Override(new(job.Service))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}
