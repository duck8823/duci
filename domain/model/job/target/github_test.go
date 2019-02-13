package target_test

import (
	"context"
	"errors"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/job/target/git/mock_git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	"testing"
)

func TestGithubPush_Prepare(t *testing.T) {
	t.Run("when success git clone", func(t *testing.T) {
		// given
		repo := &target.MockRepository{
			FullName: "duck8823/duci",
			URL:      "http://example.com",
		}
		point := &github.SimpleTargetPoint{
			Ref: "test",
			SHA: random.String(16, random.Alphanumeric),
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		mockGit := mock_git.NewMockGit(ctrl)
		mockGit.EXPECT().
			Clone(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
		container.Override(mockGit)
		defer container.Clear()

		// and
		sut := &target.GitHub{
			Repo:  repo,
			Point: point,
		}

		// when
		got, cleanup, err := sut.Prepare(context.Background())
		defer cleanup()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if len(got) == 0 {
			t.Error("must not be empty")
		}
	})

	t.Run("when failure git clone", func(t *testing.T) {
		// given
		repo := &target.MockRepository{
			FullName: "duck8823/duci",
			URL:      "http://example.com",
		}
		point := &github.SimpleTargetPoint{
			Ref: "test",
			SHA: random.String(16, random.Alphanumeric),
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		mockGit := mock_git.NewMockGit(ctrl)
		mockGit.EXPECT().
			Clone(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("test error"))
		container.Override(mockGit)
		defer container.Clear()

		// and
		sut := &target.GitHub{
			Repo:  repo,
			Point: point,
		}

		// when
		got, cleanup, err := sut.Prepare(context.Background())
		defer cleanup()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if len(got) != 0 {
			t.Errorf("must be empty, but got %+v", got)
		}
	})

	t.Run("when git have not be initialized", func(t *testing.T) {
		// given
		repo := &target.MockRepository{
			FullName: "duck8823/duci",
			URL:      "http://example.com",
		}
		point := &github.SimpleTargetPoint{
			Ref: "test",
			SHA: random.String(16, random.Alphanumeric),
		}

		// and
		container.Clear()

		// and
		sut := &target.GitHub{
			Repo:  repo,
			Point: point,
		}

		// when
		got, cleanup, err := sut.Prepare(context.Background())
		defer cleanup()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if len(got) != 0 {
			t.Errorf("must be empty, but got %+v", got)
		}
	})
}
