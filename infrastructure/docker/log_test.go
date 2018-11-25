package docker_test

import (
	"bufio"
	"bytes"
	"github.com/duck8823/duci/infrastructure/docker"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildLogger_ReadLine(t *testing.T) {
	// given
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}
	date := time.Date(2020, time.December, 4, 4, 32, 12, 3, jst)

	docker.SetNowFunc(func() time.Time {
		return date
	})
	defer docker.SetNowFunc(time.Now)

	// and
	reader := bufio.NewReader(strings.NewReader("{\"stream\":\"Hello World.\"}"))
	logger := &docker.BuildLogger{}
	logger.SetReader(reader)

	// and
	expected := &docker.LogLine{Timestamp: date, Message: []byte("Hello World.")}

	// when
	actual, err := logger.ReadLine()

	// then
	if err != nil {
		t.Errorf("error must not occur, but got %+v", err)
	}

	// and
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("must be equal: wont %+v, but got %+v", expected, actual)
	}
}

func TestRunLogger_ReadLine(t *testing.T) {
	// setup
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}
	date := time.Date(2020, time.December, 4, 4, 32, 12, 3, jst)

	docker.SetNowFunc(func() time.Time {
		return date
	})
	defer docker.SetNowFunc(time.Now)

	t.Run("with correct format", func(t *testing.T) {
		// given
		prefix := []byte{1, 0, 0, 0, 9, 9, 9, 9}
		reader := bufio.NewReader(bytes.NewReader(append(prefix, 'H', 'e', 'l', 'l', 'o')))
		logger := &docker.RunLogger{}
		logger.SetReader(reader)

		// and
		expected := &docker.LogLine{Timestamp: date, Message: []byte("Hello")}

		// when
		actual, err := logger.ReadLine()

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// and
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("must be equal: wont %+v, but got %+v", expected, actual)
		}
	})

	t.Run("with invalid format", func(t *testing.T) {
		// given
		prefix := []byte{0, 0, 0, 0, 9, 9, 9, 9}
		reader := bufio.NewReader(bytes.NewReader(append(prefix, 'H', 'e', 'l', 'l', 'o')))
		logger := &docker.RunLogger{}
		logger.SetReader(reader)

		// when
		actual, err := logger.ReadLine()

		// then
		if err == nil {
			t.Error("error must occur, but got nil")
		}

		// and
		if actual != nil {
			t.Errorf("must be equal: wont nil, but got %+v", actual)
		}
	})

	t.Run("when too short", func(t *testing.T) {
		// given
		reader := bufio.NewReader(bytes.NewReader([]byte{'H', 'e', 'l', 'l', 'o'}))
		logger := &docker.RunLogger{}
		logger.SetReader(reader)

		// and
		expected := &docker.LogLine{Timestamp: date, Message: []byte{}}

		// when
		actual, err := logger.ReadLine()

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// and
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("must be equal: wont %+v, but got %+v", expected, actual)
		}
	})
}
