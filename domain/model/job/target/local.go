package target

import (
	"github.com/duck8823/duci/domain/model/job"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// Local is target with Local directory
type Local struct {
	Path string
}

// Prepare working directory
func (l *Local) Prepare() (job.WorkDir, job.Cleanup, error) {
	tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric, random.Numeric))
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", nil, errors.WithStack(err)
	}

	if err := copyDir(tmpDir, l.Path); err != nil {
		return "", nil, errors.WithStack(err)
	}

	return job.WorkDir(tmpDir), cleanupFunc(tmpDir), nil
}

func copyDir(dstDir string, srcDir string) error {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, entry := range entries {
		dstPath := path.Join(dstDir, entry.Name())
		srcPath := path.Join(srcDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0700); err != nil {
				return errors.WithStack(err)
			}

			if err := copyDir(dstPath, srcPath); err != nil {
				return errors.WithStack(err)
			}
		} else if err := copyFile(dstPath, srcPath); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func copyFile(dstFile string, srcFile string) error {
	dst, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dst.Close()

	src, err := os.Open(srcFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func cleanupFunc(path string) job.Cleanup {
	return func() {
		if err := os.RemoveAll(path); err != nil {
			logrus.Error(err)
		}
	}
}
