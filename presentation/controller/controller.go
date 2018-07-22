package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/logger"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strings"
)

var SKIP_BUILD = errors.New("build skip")

type JobController struct {
	Runner runner.Runner
	GitHub github.Service
}

func (c *JobController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestId := uuid.New()

	// Trigger build
	githubEvent := r.Header.Get("X-GitHub-Event")
	switch githubEvent {
	case "issue_comment":
		// Read Payload
		event := &go_github.IssueCommentEvent{}
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, repo, ref, command, err := c.parseIssueComment(event)
		if err == SKIP_BUILD {
			logger.Info(requestId, "skip build")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		} else if err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go c.Runner.Run(ctx, repo, ref, command...)
	case "push":
		event := &go_github.PushEvent{}
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go c.Runner.Run(context.New("push"), event.GetRepo(), event.GetRef())
	default:
		message := fmt.Sprintf("payload event type must be issue_comment or push. but %s", githubEvent)
		logger.Error(requestId, message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}

func (c *JobController) parseIssueComment(event *go_github.IssueCommentEvent) (ctx context.Context, repo *go_github.Repository, ref string, command []string, err error) {
	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		return nil, nil, "", nil, SKIP_BUILD
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")
	command = strings.Split(phrase, " ")
	ctx = context.New(fmt.Sprintf("pr/%s", command[0]))

	pr, err := c.GitHub.GetPullRequest(ctx, event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, nil, "", nil, errors.WithStack(err)
	}

	repo = event.GetRepo()
	ref = pr.GetHead().GetRef()
	return ctx, repo, ref, command, err
}
