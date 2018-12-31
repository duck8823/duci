package runner_test

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/gommon/random"
	"os"
	"path"
	"testing"
)

func TestCreateTarball(t *testing.T) {
	t.Run("with correct directory", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		if err := os.MkdirAll(tmpDir, 0700); err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		defer os.RemoveAll(tmpDir)

		if _, err := os.Create(path.Join(tmpDir, "a")); err != nil {
			t.Fatalf("error occur: %+v", err)
		}

		// and
		want := path.Join(tmpDir, "duci.tar")

		// when
		got, err := runner.CreateTarball(job.WorkDir(tmpDir))

		// then
		if err != nil {
			t.Fatalf("error must be nil, but got %+v", err)
		}
		defer got.Close()

		// and
		if got.Name() != want {
			t.Errorf("file name: want %s, but got %s", want, got.Name())
		}
	})

	t.Run("with invalid directory", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		// when
		got, err := runner.CreateTarball(job.WorkDir(tmpDir))

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)

			got.Close()
		}
	})
}

func TestDockerfilePath(t *testing.T) {
	// where
	for _, tt := range []struct {
		name  string
		given func(t *testing.T) (workDir job.WorkDir, cleanup func())
		want  docker.Dockerfile
	}{
		{
			name: "when .duci directory not found",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(tmpDir, 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want: "./Dockerfile",
		},
		{
			name: "when .duci directory found but .duci/Dockerfile not found",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(path.Join(tmpDir, ".duci"), 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want: "./Dockerfile",
		},
		{
			name: "when .duci/Dockerfile found",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(path.Join(tmpDir, ".duci"), 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}
				if _, err := os.Create(path.Join(tmpDir, ".duci", "Dockerfile")); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want: ".duci/Dockerfile",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// given
			in, cleanup := tt.given(t)

			// when
			got := runner.DockerfilePath(in)

			// then
			if got != tt.want {
				t.Errorf("must be equal, but %+v", cmp.Diff(got, tt.want))
			}

			// cleanup
			cleanup()
		})
	}
}

func TestRuntimeOptions(t *testing.T) {
	// where
	for _, tt := range []struct {
		name    string
		given   func(t *testing.T) (workDir job.WorkDir, cleanup func())
		want    docker.RuntimeOptions
		wantErr bool
	}{
		{
			name: "when .duci/config.yml not found",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(tmpDir, 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want:    docker.RuntimeOptions{},
			wantErr: false,
		},
		{
			name: "when .duci/config.yml found",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(path.Join(tmpDir, ".duci"), 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				file, err := os.OpenFile(path.Join(tmpDir, ".duci", "config.yml"), os.O_RDWR|os.O_CREATE, 0400)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				defer file.Close()

				file.WriteString(`---
volumes:
  - hoge:fuga
`)

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want: docker.RuntimeOptions{
				Volumes: docker.Volumes{"hoge:fuga"},
			},
			wantErr: false,
		},
		{
			name: "when .duci/config.yml is directory",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(path.Join(tmpDir, ".duci", "config.yml"), 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want:    docker.RuntimeOptions{},
			wantErr: true,
		},
		{
			name: "when .duci/config.yml is invalid format",
			given: func(t *testing.T) (workDir job.WorkDir, cleanup func()) {
				t.Helper()

				tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
				if err := os.MkdirAll(path.Join(tmpDir, ".duci"), 0700); err != nil {
					t.Fatalf("error occur: %+v", err)
				}

				file, err := os.OpenFile(path.Join(tmpDir, ".duci", "config.yml"), os.O_RDWR|os.O_CREATE, 0400)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				defer file.Close()

				file.WriteString("invalid format")

				return job.WorkDir(tmpDir), func() {
					_ = os.RemoveAll(tmpDir)
				}
			},
			want:    docker.RuntimeOptions{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// given
			in, cleanup := tt.given(t)

			// when
			got, err := runner.ExportedRuntimeOptions(in)

			// then
			if tt.wantErr && err == nil {
				t.Error("error must not be nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			// and
			if !cmp.Equal(got, tt.want) {
				t.Errorf("must be equal, but %+v", cmp.Diff(got, tt.want))
			}

			// cleanup
			cleanup()
		})
	}
}
