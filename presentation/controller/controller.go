package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/docker"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/service/github"
	"github.com/duck8823/minimal-ci/service/runner"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var SKIP_BUILD = errors.New("build skip")

type jobController struct {
	runner runner.Runner
	github github.Service
}

func New() (*jobController, error) {
	name := "minimal-ci"

	githubService, err := github.New(os.Getenv("GITHUB_API_TOKEN"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dockerClient, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dockerRunner := &runner.DockerRunner{
		Name:        name,
		BaseWorkDir: path.Join(os.TempDir(), name),
		GitHub:      githubService,
		Docker:      dockerClient,
	}

	return &jobController{dockerRunner, githubService}, nil
}

func (c *jobController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		go c.runner.Run(ctx, repo, ref, command...)
	case "push":
		event := &go_github.PushEvent{}
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			logger.Errorf(requestId, "%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go c.runner.Run(context.New("push"), event.GetRepo(), event.GetRef())
	default:
		message := fmt.Sprintf("payload event type must be issue_comment or push. but %s", githubEvent)
		logger.Error(requestId, message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}

func (c *jobController) parseIssueComment(event *go_github.IssueCommentEvent) (ctx context.Context, repo *go_github.Repository, ref string, command []string, err error) {
	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
		return nil, nil, "", nil, SKIP_BUILD
	}
	phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")
	command = strings.Split(phrase, " ")
	ctx = context.New(fmt.Sprintf("pr/%s", command[0]))

	pr, err := c.github.GetPullRequest(ctx, event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, nil, "", nil, errors.WithStack(err)
	}

	repo = event.GetRepo()
	ref = pr.GetHead().GetRef()
	return
}
