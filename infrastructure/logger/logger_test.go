package logger_test

import (
	"github.com/duck8823/minimal-ci/infrastructure/logger"
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
	InitLogger(t)

	logger.Debug(uuid.UUID{}, "Hello World.")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[36;1m[DEBUG]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestDebugf(t *testing.T) {
	InitLogger(t)

	logger.Debugf(uuid.UUID{}, "Hello %s.", "World")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[36;1m[DEBUG]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestInfo(t *testing.T) {
	InitLogger(t)

	logger.Info(uuid.UUID{}, "Hello World.")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[1m[INFO]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestInfof(t *testing.T) {
	InitLogger(t)

	logger.Infof(uuid.UUID{}, "Hello %s.", "World")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[1m[INFO]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestError(t *testing.T) {
	InitLogger(t)

	logger.Error(uuid.UUID{}, "Hello World.")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[41;1m[ERROR]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func TestErrorf(t *testing.T) {
	InitLogger(t)

	logger.Errorf(uuid.UUID{}, "Hello %s.", "World")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := "[00000000-0000-0000-0000-000000000000]  \033[41;1m[ERROR]\033[0m Hello World."
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func InitLogger(t *testing.T) {
	t.Helper()

	reader, writer, _ = os.Pipe()

	logger.Writer = writer
}

func ReadLog(t *testing.T) string {
	t.Helper()

	writer.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error()
	}

	return string(bytes)
}
