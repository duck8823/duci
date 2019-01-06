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
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
)

// ErrSkipBuild represents error of skip build
var ErrSkipBuild = errors.New("Skip build")

type handler struct {
	executor executor.Executor
}

// NewHandler returns a implement of webhook handler
func NewHandler() (http.Handler, error) {
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

	targetURL := targetURL(r)
	targetURL.Path = fmt.Sprintf("/logs/%s", reqID.ToSlice())
	ctx := application.ContextWithJob(context.Background(), &application.BuildJob{
		ID: reqID,
		TargetSource: &github.TargetSource{
			Repository: event.GetRepo(),
			Ref:        event.GetRef(),
			SHA:        plumbing.NewHash(event.GetHeadCommit().GetID()),
		},
		TaskName:  fmt.Sprintf("%s/push", application.Name),
		TargetURL: targetURL,
	})

	tgt := &target.GitHub{
		Repo:  event.GetRepo(),
		Point: event,
	}

	go func() {
		if err := h.executor.Execute(ctx, tgt); err != nil {
			logrus.Error(err)
		}
	}()

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
		if _, err := w.Write([]byte("{\"message\":\"skip build\"}")); err != nil {
			logrus.Error(err)
		}
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
	if err == ErrSkipBuild {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("{\"message\":\"skip build\"}")); err != nil {
			logrus.Error(err)
		}
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	targetURL := targetURL(r)
	targetURL.Path = fmt.Sprintf("/logs/%s", reqID.ToSlice())
	ctx := application.ContextWithJob(context.Background(), &application.BuildJob{
		ID: reqID,
		TargetSource: &github.TargetSource{
			Repository: event.GetRepo(),
			Ref:        pnt.GetRef(),
			SHA:        plumbing.NewHash(pnt.GetHead()),
		},
		TaskName:  fmt.Sprintf("%s/pr/%s", application.Name, phrase.Command().Slice()[0]),
		TargetURL: targetURL,
	})

	tgt := &target.GitHub{
		Repo:  event.GetRepo(),
		Point: pnt,
	}

	go func() {
		if err := h.executor.Execute(ctx, tgt, phrase.Command()...); err != nil {
			logrus.Error(err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}
