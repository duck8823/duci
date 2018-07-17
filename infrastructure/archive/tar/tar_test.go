package tar_test

import (
	archiveTar "archive/tar"
	"fmt"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	testDir := createTestDir(t)
	archiveDir := path.Join(testDir, "archive")

	createFile(t, path.Join(archiveDir, "file"), "this is file.")
	createFile(t, path.Join(archiveDir, "dir", "file"), "this is file in the dir.")

	if err := os.MkdirAll(path.Join(archiveDir, "empty"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	output := path.Join(testDir, "output.tar")
	tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer tarFile.Close()

	if err := tar.Create(archiveDir, tarFile); err != nil {
		t.Fatalf("%+v", err)
	}

	file, err := os.Open(output)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	actual := readTarArchive(t, file)

	expect := Files{
		{
			Name:    "dir/file",
			Content: "this is file in the dir.",
		},
		{
			Name:    "file",
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

func readTarArchive(t *testing.T, reader io.Reader) Files {
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

func createTestDir(t *testing.T) string {
	t.Helper()

	tempDir := path.Join(os.TempDir(), fmt.Sprintf("duci_test_%v", time.Now().Unix()))
	if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	return tempDir
}

func createFile(t *testing.T, name string, content string) {
	t.Helper()

	paths := strings.Split(name, "/")
	dir := strings.Join(paths[:len(paths)-1], "/")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("%+v", err)
	}
}
