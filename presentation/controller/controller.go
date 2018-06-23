package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duck8823/minimal-ci/service/runner"
	"github.com/google/go-github/github"
	"github.com/google/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

type jobController struct {
	runner *runner.Runner
}

func New() (*jobController, error) {
	jobRunner, err := runner.NewWithEnv()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &jobController{jobRunner}, nil
}

func (c *jobController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read Payload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := &github.IssueCommentEvent{}
	if err := json.Unmarshal(body, event); err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trigger build
	githubEvent := r.Header.Get("X-GitHub-Event")
	if githubEvent != "issue_comment" {
		message := fmt.Sprintf("payload event type must be issue_comment. but %s", githubEvent)
		logger.Error(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not build."))
		return
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")

	if err := c.runner.Run(context.Background(), event.GetRepo(), "master", phrase); err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}
