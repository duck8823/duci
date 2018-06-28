package logger_test

import (
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"github.com/google/uuid"
)

var (
	reader io.ReadCloser
	writer io.WriteCloser
	regex  = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}.\\d{3}")
)

func TestDebug(t *testing.T) {
	InitLogger(t)

	logger.Debug(uuid.UUID{}, "Hello World.")

	log := ReadLog(t)
	if !regex.MatchString(log) {
		t.Fatalf("invalid format. %+v", log)
	}
	actual := strings.TrimRight(regex.ReplaceAllString(log, ""), "\n")

	expected := fmt.Sprintf("%s [DEBUG]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", logging.ColorSeq(logging.ColorCyan))
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

	expected := fmt.Sprintf("%s [DEBUG]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", logging.ColorSeq(logging.ColorCyan))
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

	expected := fmt.Sprintf("%s [INFO]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", "")
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

	expected := fmt.Sprintf("%s [INFO]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", "")
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

	expected := fmt.Sprintf("%s [ERROR]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", logging.ColorSeq(logging.ColorRed))
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

	expected := fmt.Sprintf("%s [ERROR]\033[0m [00000000-0000-0000-0000-000000000000] Hello World.", logging.ColorSeq(logging.ColorRed))
	if actual != expected {
		t.Errorf("wrong log. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}
}

func InitLogger(t *testing.T) {
	t.Helper()

	reader, writer, _ = os.Pipe()

	logger.Init(writer, logging.DEBUG)
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
