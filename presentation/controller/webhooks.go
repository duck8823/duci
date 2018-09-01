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
	deliveryId := go_github.DeliveryID(r)
	requestId, err := uuid.Parse(deliveryId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: invalid request header `X-GitHub-Delivery`: %+v", deliveryId), http.StatusBadRequest)
		return
	}

	runtimeUrl := &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.URL.Path,
	}
	if r.URL.Scheme != "" {
		runtimeUrl.Scheme = r.URL.Scheme
	}

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

		ctx, repo, head, command, err := c.parseIssueComment(event, requestId, runtimeUrl)
		if err == SkipBuild {
			logger.Info(requestId, "skip build")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		} else if err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ref := fmt.Sprintf("refs/heads/%s", head.GetRef())
		go c.Runner.Run(ctx, repo, ref, plumbing.NewHash(head.GetSHA()), command...)
	case "push":
		event := &go_github.PushEvent{}
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sha := event.GetHeadCommit().GetID()
		if len(sha) == 0 {
			logger.Info(requestId, "skip build")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("skip build"))
			return
		}

		taskName := fmt.Sprintf("%s/push", application.Name)
		ctx := context.New(taskName, requestId, runtimeUrl)
		go c.Runner.Run(ctx, event.GetRepo(), event.GetRef(), plumbing.NewHash(sha))
	default:
		message := fmt.Sprintf("payload event type must be issue_comment or push. but %s", githubEvent)
		logger.Error(requestId, message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}

func (c *WebhooksController) parseIssueComment(
	event *go_github.IssueCommentEvent,
	requestId uuid.UUID,
	url *url.URL,
) (ctx context.Context, repo *go_github.Repository, head *go_github.PullRequestBranch, command []string, err error) {

	if !isValidAction(event.Action) {
		return nil, nil, nil, nil, SkipBuild
	}

	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		return nil, nil, nil, nil, SkipBuild
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")
	command = strings.Split(phrase, " ")
	ctx = context.New(fmt.Sprintf("%s/pr/%s", application.Name, command[0]), requestId, url)

	pr, err := c.GitHub.GetPullRequest(ctx, event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, nil, nil, nil, errors.WithStack(err)
	}

	repo = event.GetRepo()
	head = pr.GetHead()
	return ctx, repo, head, command, err
}

func isValidAction(action *string) bool {
	if action == nil {
		return false
	}
	return *action == "created" || *action == "edited"
}
