package controller

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/minimal-ci/service/runner/mock_runner"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJobController_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("with correct payload", func(t *testing.T) {
		mock := mock_runner.NewMockRunner(ctrl)
		mock.EXPECT().ConvertPullRequestToRef(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("master", nil)
		mock.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		handler := &jobController{runner: mock}

		s := httptest.NewServer(handler)
		defer s.Close()

		body := CreateBody(t, "ci test")

		req := httptest.NewRequest("POST", "/", body)
		req.Header.Set("X-GitHub-Event", "issue_comment")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		actual := rec.Code
		expected := 200
		if actual != expected {
			t.Errorf("status must equal %+v, but got %+v", expected, actual)
		}
	})

	t.Run("when runner returns error", func(t *testing.T) {
		mock := mock_runner.NewMockRunner(ctrl)
		mock.EXPECT().ConvertPullRequestToRef(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("", errors.New("error"))

		handler := &jobController{runner: mock}

		s := httptest.NewServer(handler)
		defer s.Close()

		body := CreateBody(t, "ci test")

		req := httptest.NewRequest("POST", "/", body)
		req.Header.Set("X-GitHub-Event", "issue_comment")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		actual := rec.Code
		expected := 500
		if actual != expected {
			t.Errorf("status must equal %+v, but got %+v", expected, actual)
		}
	})

	t.Run("must not call RunWithPullRequest", func(t *testing.T) {
		mock := mock_runner.NewMockRunner(ctrl)
		mock.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		handler := &jobController{runner: mock}

		s := httptest.NewServer(handler)
		defer s.Close()

		t.Run("with invalid header", func(t *testing.T) {
			body := CreateBody(t, "ci test")

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "hogefuga")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 500
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})

		t.Run("without comment started ci", func(t *testing.T) {
			body := CreateBody(t, "test")

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "issue_comment")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 200
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})

		t.Run("with invalid body", func(t *testing.T) {
			body := strings.NewReader("Invalid JSON format.")

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "issue_comment")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 500
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})
	})
}

func CreateBody(t *testing.T, comment string) io.Reader {
	t.Helper()

	number := 1
	event := &github.IssueCommentEvent{
		Repo: &github.Repository{},
		Issue: &github.Issue{
			Number: &number,
		},
		Comment: &github.IssueComment{
			Body: &comment,
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	return bytes.NewReader(payload)
}
