package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/logger"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var SkipBuild = errors.New("build skip")

type WebhooksController struct {
	Runner runner.Runner
	GitHub github.Service
}

func (c *WebhooksController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID, err := requestID(r)
	if err != nil {
		logger.Error(requestID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Trigger build
	githubEvent := r.Header.Get("X-GitHub-Event")
	switch githubEvent {
	case "issue_comment":
		c.runWithIssueCommentEvent(requestID, w, r)
	case "push":
		c.runWithPushEvent(requestID, w, r)
	default:
		message := fmt.Sprintf("payload event type must be issue_comment or push. but %s", githubEvent)
		logger.Error(requestID, message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}

func (c *WebhooksController) runWithIssueCommentEvent(
	requestID uuid.UUID,
	w http.ResponseWriter,
	r *http.Request,
) {
	// Read Payload
	event := &go_github.IssueCommentEvent{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		logger.Errorf(requestID, "%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, repo, head, command, err := c.parseIssueComment(event, requestID, runtimeURL(r))
	if err == SkipBuild {
		logger.Info(requestID, "skip build")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		logger.Errorf(requestID, "%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ref := fmt.Sprintf("refs/heads/%s", head.GetRef())
	go c.Runner.Run(ctx, repo, ref, plumbing.NewHash(head.GetSHA()), command...)
}

func (c *WebhooksController) parseIssueComment(
	event *go_github.IssueCommentEvent,
	requestID uuid.UUID,
	url *url.URL,
) (ctx context.Context, repo *go_github.Repository, head *go_github.PullRequestBranch, command []string, err error) {
	phrase, err := phrase(event)
	if err != nil {
		return nil, nil, nil, nil, SkipBuild
	}
	command = strings.Split(phrase, " ")
	ctx = context.New(fmt.Sprintf("%s/pr/%s", application.Name, command[0]), requestID, url)

	pr, err := c.GitHub.GetPullRequest(ctx, event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, nil, nil, nil, errors.WithStack(err)
	}

	repo = event.GetRepo()
	head = pr.GetHead()
	return ctx, repo, head, command, err
}

func (c *WebhooksController) runWithPushEvent(requestID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	event := &go_github.PushEvent{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		logger.Errorf(requestID, "%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sha := event.GetHeadCommit().GetID()
	if len(sha) == 0 {
		logger.Info(requestID, "skip build: could not get head commit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("skip build"))
		return
	}

	taskName := fmt.Sprintf("%s/push", application.Name)
	ctx := context.New(taskName, requestID, runtimeURL(r))
	go c.Runner.Run(ctx, event.GetRepo(), event.GetRef(), plumbing.NewHash(sha))
}

func requestID(r *http.Request) (uuid.UUID, error) {
	deliveryID := go_github.DeliveryID(r)
	requestID, err := uuid.Parse(deliveryID)
	if err != nil {
		msg := fmt.Sprintf("Error: invalid request header `X-GitHub-Delivery`: %+v", deliveryID)
		return uuid.New(), errors.Wrap(err, msg)
	}
	return requestID, nil
}

func runtimeURL(r *http.Request) *url.URL {
	runtimeURL := &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.URL.Path,
	}
	if r.URL.Scheme != "" {
		runtimeURL.Scheme = r.URL.Scheme
	}
	return runtimeURL
}

func phrase(event *go_github.IssueCommentEvent) (string, error) {
	if !isValidAction(event.Action) {
		return "", SkipBuild
	}

	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		return "", SkipBuild
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")
	return phrase, nil
}

func isValidAction(action *string) bool {
	if action == nil {
		return false
	}
	return *action == "created" || *action == "edited"
}
