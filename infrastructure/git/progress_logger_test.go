package git_test

import (
	"github.com/duck8823/duci/infrastructure/git"
	"github.com/duck8823/duci/infrastructure/logger"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

var regex = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}`)

func TestProgressLogger_Write(t *testing.T) {
	// setup
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	defer reader.Close()

	logger.Writer = writer

	// given
	progress := &git.ProgressLogger{}

	// when
	progress.Write([]byte("hoge\rfuga"))
	writer.Close()

	actual := readLogLine(t, reader)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[1m[INFO]\033[0m hoge"

	// then
	if actual != expected {
		t.Errorf("must remove CR flag or later. wont: %+v, but got: %+v", expected, actual)
	}
}

func readLogLine(t *testing.T, reader io.Reader) string {
	t.Helper()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	log := string(bytes)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}

	return strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")
}
