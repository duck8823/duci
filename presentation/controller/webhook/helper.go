package webhook

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job/target/github"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

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

func targetPoint(event *go_github.IssueCommentEvent) (github.TargetPoint, error) {
	gh, err := github.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pr, err := gh.GetPullRequest(context.Background(), event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &github.SimpleTargetPoint{
		Ref: pr.GetHead().GetRef(),
		SHA: pr.GetHead().GetSHA(),
	}, nil
}

func isValidAction(action *string) bool {
	if action == nil {
		return false
	}
	return *action == "created" || *action == "edited"
}
