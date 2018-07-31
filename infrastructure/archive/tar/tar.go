package tar

import (
	"archive/tar"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	)

func Create(dir string, output io.Writer) error {
	writer := tar.NewWriter(output)
	defer writer.Close()

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return errors.WithStack(err)
		}

		now := clock.Now()
		header := &tar.Header{
			Name:       strings.Replace(file.Name(), dir+"/", "", -1),
			Mode:       0600,
			ModTime:    now,
			AccessTime: now,
			ChangeTime: now,
			Size:       info.Size(),
		}
		if err := writer.WriteHeader(header); err != nil {
			return errors.WithStack(err)
		}
		if _, err := writer.Write(data); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
