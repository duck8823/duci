package logger_test

import (
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	reader io.ReadCloser
	writer io.WriteCloser
)

func TestDebug(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Debug(uuid.UUID{}, "Hello World.")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[36;1m[DEBUG]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func TestDebugf(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Debugf(uuid.UUID{}, "Hello %s.", "World")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[36;1m[DEBUG]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func TestInfo(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Info(uuid.UUID{}, "Hello World.")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[1m[INFO]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func TestInfof(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Infof(uuid.UUID{}, "Hello %s.", "World")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[1m[INFO]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func TestError(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Error(uuid.UUID{}, "Hello World.")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[41;1m[ERROR]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func TestErrorf(t *testing.T) {
	// setup
	initLogger(t)

	// and
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	clock.Now = func() time.Time {
		return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
	}

	// when
	logger.Errorf(uuid.UUID{}, "Hello %s.", "World")

	actual := readLog(t)
	expected := "[00000000-0000-0000-0000-000000000000] 1987-03-27 19:19:00.000 \033[41;1m[ERROR]\033[0m Hello World."

	// then
	if actual != expected {
		t.Errorf("wrong logstore. wont: \"%+v\", got: \"%+v\"", expected, actual)
	}

	// cleanup
	clock.Adjust()
}

func initLogger(t *testing.T) {
	t.Helper()

	reader, writer, _ = os.Pipe()

	logger.Writer = writer
}

func readLog(t *testing.T) string {
	t.Helper()

	writer.Close()
	log, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error()
	}

	return strings.TrimRight(string(log), "\n")
}
