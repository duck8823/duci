package logger_test

import (
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	reader io.ReadCloser
	writer io.WriteCloser
	regex  = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}`)
)

func TestDebug(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Debug(uuid.UUID{}, "Hello World.")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[36;1m[DEBUG]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestDebugf(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Debugf(uuid.UUID{}, "Hello %s.", "World")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[36;1m[DEBUG]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestInfo(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Info(uuid.UUID{}, "Hello World.")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[1m[INFO]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestInfof(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Infof(uuid.UUID{}, "Hello %s.", "World")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[1m[INFO]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestError(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Error(uuid.UUID{}, "Hello World.")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[41;1m[ERROR]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestErrorf(t *testing.T) {
	// setup
	initLogger(t)

	// when
	logger.Errorf(uuid.UUID{}, "Hello %s.", "World")

	actual := readLogTrimmedTime(t)
	expected := "[00000000-0000-0000-0000-000000000000]  \033[41;1m[ERROR]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func initLogger(t *testing.T) {
	t.Helper()

	reader, writer, _ = os.Pipe()

	logger.Writer = writer
}

func readLogTrimmedTime(t *testing.T) string {
	t.Helper()

	writer.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error()
	}

	log := string(bytes)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	return strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")
}
