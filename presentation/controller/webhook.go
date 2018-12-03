package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/job/target/github"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/url"
)

type WebhookHandler struct {
	Executor executor.Executor
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "push":
		h.PushEvent(w, r)
	default:
		msg := fmt.Sprintf("payload event type must be push. but %s", event)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) PushEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Executor.Execute(ctx, &target.GitHubPush{
		Repo:  event.GetRepo(),
		Point: event,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reqID(r *http.Request) (uuid.UUID, error) {
	deliveryID := go_github.DeliveryID(r)
	requestID, err := uuid.Parse(deliveryID)
	if err != nil {
		msg := fmt.Sprintf("Error: invalid request header `X-GitHub-Delivery`: %+v", deliveryID)
		return uuid.New(), errors.Wrap(err, msg)
	}
	return requestID, nil
}

func targetURL(r *http.Request) *url.URL {
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
