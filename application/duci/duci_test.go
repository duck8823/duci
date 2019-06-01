package duci_test

import (
	"context"
	"errors"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/duci"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/application/service/job/mock_job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/domain/model/job/target/github/mock_github"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/container"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"net/url"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("when there are instances in container", func(t *testing.T) {
		// given
		container.Override(new(jobService.Service))
		container.Override(new(github.GitHub))
		defer container.Clear()

		// when
		got, err := duci.New()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if got == nil {
			t.Errorf("duci must not be nil")
		}
	})

	t.Run("when instance not enough in container", func(t *testing.T) {
		// where
		for _, tt := range []struct {
			name string
			in   []interface{}
		}{
			{
				name: "with only job_service.Service instance",
				in:   []interface{}{new(jobService.Service)},
			},
			{
				name: "with only github.GitHub instance",
				in:   []interface{}{new(github.GitHub)},
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// given
				container.Clear()

				for _, ins := range tt.in {
					container.Override(ins)
				}
				defer container.Clear()

				// when
				got, err := duci.New()

				// then
				if err == nil {
					t.Error("error must not be nil")
				}

				// and
				if got != nil {
					t.Errorf("duci must be nil, but got %+v", got)
				}
			})
		}
	})
}

func TestDuci_Init(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Start(gomock.Eq(buildJob.ID)).
			Times(1).
			Return(nil)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Eq(ctx), gomock.Eq(github.CommitStatus{
				TargetSource: buildJob.TargetSource,
				State:        github.PENDING,
				Description:  "queued",
				Context:      buildJob.TaskName,
				TargetURL:    buildJob.TargetURL,
			})).
			Times(1).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.Init(ctx)

		// then
		ctrl.Finish()
	})

	t.Run("when invalid build job value", func(t *testing.T) {
		// given
		ctx := context.WithValue(context.Background(), duci.String("duci_job"), "invalid value")

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Start(gomock.Any()).
			Times(0)
		service.EXPECT().
			Append(gomock.Any(), gomock.Any()).
			Times(0)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.Init(ctx)

		// then
		ctrl.Finish()
	})

	t.Run("when failed to job_service.Service#Start", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Start(gomock.Any()).
			Times(1).
			Return(errors.New("test error"))
		service.EXPECT().
			Append(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.Init(ctx)

		// then
		ctrl.Finish()
	})
}

func TestDuci_Start(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Start(gomock.Eq(buildJob.ID)).
			Times(0)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Eq(ctx), gomock.Eq(github.CommitStatus{
				TargetSource: buildJob.TargetSource,
				State:        github.PENDING,
				Description:  "running",
				Context:      buildJob.TaskName,
				TargetURL:    buildJob.TargetURL,
			})).
			Times(1).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.Start(ctx)

		// then
		ctrl.Finish()
	})

	t.Run("when invalid build job value", func(t *testing.T) {
		// given
		ctx := context.WithValue(context.Background(), duci.String("duci_job"), "invalid value")

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Start(gomock.Any()).
			Times(0)
		service.EXPECT().
			Append(gomock.Any(), gomock.Any()).
			Times(0)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.Start(ctx)

		// then
		ctrl.Finish()
	})
}

func TestDuci_AppendLog(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)

		log := &duci.MockLog{Msgs: []string{"Hello", "World"}}

		// and
		ctrl := gomock.NewController(t)
		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Append(gomock.Eq(buildJob.ID), gomock.Any()).
			Times(len(log.Msgs)).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()

		// when
		sut.AppendLog(ctx, log)

		// then
		ctrl.Finish()
	})

	t.Run("when invalid build job value", func(t *testing.T) {
		// given
		ctx := context.WithValue(context.Background(), duci.String("duci_job"), "invalid value")
		log := &duci.MockLog{Msgs: []string{"Hello", "World"}}

		// and
		ctrl := gomock.NewController(t)
		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Append(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()

		// when
		sut.AppendLog(ctx, log)

		// then
		ctrl.Finish()
	})
}

func TestDuci_End(t *testing.T) {
	t.Run("when error is nil", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		buildJob.BeginAt(time.Unix(0, 0))
		ctx := application.ContextWithJob(context.Background(), buildJob)
		var err error = nil

		// and
		defer duci.SetNowFunc(func() time.Time {
			return time.Unix(302, 1)
		})()

		// and
		want := github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.SUCCESS,
			Description:  "success in 5min",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Finish(gomock.Any()).
			Times(1).
			Return(nil)
		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Eq(ctx), gomock.Eq(want)).
			Times(1).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.End(ctx, err)

		// then
		ctrl.Finish()
	})

	t.Run("when error is runner.Failure", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		buildJob.BeginAt(time.Unix(0, 0))
		ctx := application.ContextWithJob(context.Background(), buildJob)
		err := runner.ErrFailure

		// and
		defer duci.SetNowFunc(func() time.Time {
			return time.Unix(49, 1)
		})()

		// and
		want := github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.FAILURE,
			Description:  "failure in 49sec",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Finish(gomock.Any()).
			Times(1).
			Return(nil)
		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Eq(ctx), gomock.Eq(want)).
			Times(1).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.End(ctx, err)

		// then
		ctrl.Finish()
	})

	t.Run("when error is not nil", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)
		err := errors.New("test error")

		// and
		want := github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.ERROR,
			Description:  github.Description("error: test error"),
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Finish(gomock.Any()).
			Times(1).
			Return(nil)
		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Eq(ctx), gomock.Eq(want)).
			Times(1).
			Return(nil)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.End(ctx, err)

		// then
		ctrl.Finish()
	})

	t.Run("when invalid build job value", func(t *testing.T) {
		// given
		ctx := context.WithValue(context.Background(), duci.String("duci_job"), "invalid value")

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Finish(gomock.Any()).
			Times(0)
		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.End(ctx, nil)

		// then
		ctrl.Finish()
	})

	t.Run("when failed to job_service.Service#Finish", func(t *testing.T) {
		// given
		buildJob := &application.BuildJob{
			ID:           job.ID(uuid.New()),
			TargetSource: &github.TargetSource{},
			TaskName:     "task/name",
			TargetURL:    duci.URLMust(url.Parse("http://example.com")),
		}
		ctx := application.ContextWithJob(context.Background(), buildJob)

		// and
		ctrl := gomock.NewController(t)

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			Finish(gomock.Any()).
			Times(1).
			Return(errors.New("test error"))
		service.EXPECT().
			Append(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		hub := mock_github.NewMockGitHub(ctrl)
		hub.EXPECT().
			CreateCommitStatus(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &duci.Duci{}
		defer sut.SetJobService(service)()
		defer sut.SetGitHub(hub)()

		// when
		sut.End(ctx, nil)

		// then
		ctrl.Finish()
	})
}
