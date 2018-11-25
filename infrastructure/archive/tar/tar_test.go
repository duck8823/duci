package tar_test

import (
	archiveTar "archive/tar"
	"fmt"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("with correct target", func(t *testing.T) {
		// setup
		testDir := createTestDir(t)
		defer os.RemoveAll(testDir)

		// given
		archiveDir := filepath.Join(testDir, "archive")

		createFile(t, filepath.Join(archiveDir, "file"), "this is file.", 0400)
		createFile(t, filepath.Join(archiveDir, "dir", "file"), "this is file in the dir.", 0400)

		if err := os.MkdirAll(filepath.Join(archiveDir, "empty"), 0700); err != nil {
			t.Fatalf("%+v", err)
		}

		output := filepath.Join(testDir, "output.tar")
		tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer tarFile.Close()

		// and
		expected := Files{
			{
				Name:    filepath.Join("dir", "file"),
				Content: "this is file in the dir.",
			},
			{
				Name:    "file",
				Content: "this is file.",
			},
		}

		// when
		if err := tar.Create(archiveDir, tarFile); err != nil {
			t.Fatalf("%+v", err)
		}
		actual := readTarArchive(t, output)

		// then
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("wrong tar contents.\nactual: %+v\nwont: %+v", actual, expected)
		}
	})

	t.Run("with wrong directory path", func(t *testing.T) {
		// setup
		testDir := createTestDir(t)
		defer os.RemoveAll(testDir)

		// given
		output := filepath.Join(testDir, "output.tar")
		tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer tarFile.Close()

		// expect
		if err := tar.Create("/path/to/wrong/dir", tarFile); err == nil {
			t.Error("error must occur")
		}
	})

	t.Run("with closed output", func(t *testing.T) {
		// setup
		testDir := createTestDir(t)
		defer os.RemoveAll(testDir)

		// given
		archiveDir := filepath.Join(testDir, "archive")

		createFile(t, filepath.Join(archiveDir, "file"), "this is file.", 0400)

		output := filepath.Join(testDir, "output.tar")
		tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		tarFile.Close()

		// expect
		if err := tar.Create(archiveDir, tarFile); err == nil {
			t.Error("error must occur")
		}
	})

	t.Run("with wrong permission in target", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skip if root user")
		}

		// setup
		testDir := createTestDir(t)
		defer os.RemoveAll(testDir)

		// given
		archiveDir := filepath.Join(testDir, "archive")

		createFile(t, filepath.Join(archiveDir, "file"), "this is file.", 0000)

		output := filepath.Join(testDir, "output.tar")
		tarFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0400)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer tarFile.Close()

		// expect
		if err := tar.Create(archiveDir, tarFile); err == nil {
			t.Error("error must occur")
		}
	})
}

type Files []struct {
	Name    string
	Content string
}

func readTarArchive(t *testing.T, output string) Files {
	file, err := os.Open(output)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var files Files

	tarReader := archiveTar.NewReader(file)
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

	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("duci_test_%v", time.Now().Unix()))
	if err := os.MkdirAll(filepath.Join(tempDir, "dir"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	return tempDir
}

func createFile(t *testing.T, name string, content string, perm os.FileMode) {
	t.Helper()

	paths := strings.Split(name, string(os.PathSeparator))
	dir := strings.Join(paths[:len(paths)-1], string(os.PathSeparator))
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("%+v", err)
	}
}
