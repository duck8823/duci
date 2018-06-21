package tar_test

import (
	archiveTar "archive/tar"
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	tempDir := CreateTestDir(t)

	if err := os.Mkdir(path.Join(tempDir, "empty"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}
	CreateFile(t, path.Join(tempDir, "file"), "this is file.")
	CreateFile(t, path.Join(tempDir, "dir", "file"), "this is file in the dir.")

	output := path.Join(os.TempDir(), "output.tar")
	tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer tarFile.Close()

	if err := tar.Create(tempDir, tarFile); err != nil {
		t.Fatalf("%+v", err)
	}

	file, err := os.Open(output)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	actual := ReadTarArchive(t, file)

	expect := Files{
		{
			Name:    "/dir/file",
			Content: "this is file in the dir.",
		},
		{
			Name:    "/file",
			Content: "this is file.",
		},
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("wrong tar contents.\nactual: %+v\nwont: %+v", actual, expect)
	}
}

type Files []struct {
	Name    string
	Content string
}

func ReadTarArchive(t *testing.T, reader io.Reader) Files {
	var files Files

	tarReader := archiveTar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		content, err := ioutil.ReadAll(tarReader)
		if err != nil {
			t.Fatalf("%+v, %+v", err, header)
		}

		files = append(files, struct {
			Name    string
			Content string
		}{Name: header.Name, Content: string(content)})
	}

	return files
}

func CreateTestDir(t *testing.T) string {
	t.Helper()

	tempDir := path.Join(os.TempDir(), fmt.Sprintf("minimal-ci_test_%v", time.Now().Unix()))
	if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	return tempDir
}

func CreateFile(t *testing.T, name string, content string) {
	t.Helper()

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("%+v", err)
	}
}
