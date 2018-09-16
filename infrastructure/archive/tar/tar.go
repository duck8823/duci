package tar

import (
	"archive/tar"
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
		content, err := newContent(path, dir)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := content.write(writer); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type content struct {
	Header *tar.Header
	Data   []byte
}

func newContent(path string, dir string) (*content, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	header := &tar.Header{
		Name: strings.Replace(file.Name(), dir+"/", "", -1),
		Mode: 0600,
		Size: int64(len(data)),
	}
	return &content{Header: header, Data: data}, nil
}

func (c *content) write(w *tar.Writer) error {
	if err := w.WriteHeader(c.Header); err != nil {
		return errors.WithStack(err)
	}
	if _, err := w.Write(c.Data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
