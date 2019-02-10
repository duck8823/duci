package webhook

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/job/target/github"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

func reqID(r *http.Request) (job.ID, error) {
	deliveryID := go_github.DeliveryID(r)
	requestID, err := uuid.Parse(deliveryID)
	if err != nil {
		msg := fmt.Sprintf("Error: invalid request header `X-GitHub-Delivery`: %+v", deliveryID)
		return job.ID{}, errors.Wrap(err, msg)
	}
	return job.ID(requestID), nil
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

func targetGitHub(event *go_github.IssueCommentEvent) (*target.GitHub, error) {
	gh, err := github.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pr, err := gh.GetPullRequest(context.Background(), event.GetRepo(), event.GetIssue().GetNumber())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &target.GitHub{
		Repo: pr.GetHead().GetRepo(),
		Point: &github.SimpleTargetPoint{
			Ref: fmt.Sprintf("refs/heads/%s", pr.GetHead().GetRef()),
			SHA: pr.GetHead().GetSHA(),
		},
	}, nil
}

func isValidAction(action *string) bool {
	if action == nil {
		return false
	}
	return *action == "created" || *action == "edited"
}
