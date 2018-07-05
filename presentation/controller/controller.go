package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/service/runner"
	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
)

type jobController struct {
	runner runner.Runner
}

func New() (*jobController, error) {
	jobRunner, err := runner.NewWithEnv()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &jobController{jobRunner}, nil
}

func (c *jobController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestId := uuid.New()

	// Read Payload
	event := &github.IssueCommentEvent{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trigger build
	githubEvent := r.Header.Get("X-GitHub-Event")
	if githubEvent != "issue_comment" {
		message := fmt.Sprintf("payload event type must be issue_comment. but %s", githubEvent)
		logger.Error(requestId, message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		logger.Info(requestId, "no build.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not build."))
		return
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")

	ref, err := c.runner.ConvertPullRequestToRef(context.New(), event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		logger.Errorf(requestId, "%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	go c.runner.Run(context.New(), event.GetRepo(), ref, phrase)

	// Response
	w.WriteHeader(http.StatusOK)
}
