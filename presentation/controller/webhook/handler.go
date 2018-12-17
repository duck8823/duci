package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/duci"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/job/target/github"
	go_github "github.com/google/go-github/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
)

// SkipBuild represents error of skip build
var SkipBuild = errors.New("Skip build")

type handler struct {
	executor executor.Executor
}

// NewHandler returns a implement of webhook handler
func NewHandler() (*handler, error) {
	executor, err := duci.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &handler{executor: executor}, nil
}

// ServeHTTP receives github event
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "push":
		h.PushEvent(w, r)
	case "issue_comment":
		h.IssueCommentEvent(w, r)
	default:
		msg := fmt.Sprintf("payload event type must be push or issue_comment. but %s", event)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

// PushEvent receives github push event
func (h *handler) PushEvent(w http.ResponseWriter, r *http.Request) {
	event := &go_github.PushEvent{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqID, err := reqID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := application.ContextWithJob(context.Background(), &application.BuildJob{
		ID: reqID,
		TargetSource: &github.TargetSource{
			Repository: event.GetRepo(),
			Ref:        event.GetRef(),
			SHA:        plumbing.NewHash(event.GetHead()),
		},
		TaskName:  fmt.Sprintf("%s/push", application.Name),
		TargetURL: targetURL(r),
	})

	tgt := &target.GitHubPush{
		Repo:  event.GetRepo(),
		Point: event,
	}

	if err := h.executor.Execute(ctx, tgt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// IssueCommentEvent receives github issue comment event
func (h *handler) IssueCommentEvent(w http.ResponseWriter, r *http.Request) {
	event := &go_github.IssueCommentEvent{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidAction(event.Action) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\":\"skip build\"}"))
		return
	}

	reqID, err := reqID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pnt, err := targetPoint(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	phrase, err := extractBuildPhrase(event.GetComment().GetBody())
	if err == SkipBuild {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\":\"skip build\"}"))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := application.ContextWithJob(context.Background(), &application.BuildJob{
		ID: reqID,
		TargetSource: &github.TargetSource{
			Repository: event.GetRepo(),
			Ref:        pnt.GetRef(),
			SHA:        plumbing.NewHash(pnt.GetHead()),
		},
		TaskName:  fmt.Sprintf("%s/pr/%s", application.Name, phrase.Command().Slice()[0]),
		TargetURL: targetURL(r),
	})

	tgt := &target.GitHubPush{
		Repo:  event.GetRepo(),
		Point: pnt,
	}

	go h.executor.Execute(ctx, tgt, phrase.Command()...)

	w.WriteHeader(http.StatusOK)
}
