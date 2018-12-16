package target_test

import (
	"errors"
	"github.com/duck8823/duci/domain/internal/container"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/job/target/git/mock_git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/gommon/random"
	"testing"
)

func TestNewGithubPush(t *testing.T) {
	t.Run("when git found", func(t *testing.T) {
		// given
		wantGit := &target.StubGit{}

		// and
		container.Override(wantGit)
		defer container.Clear()

		// and
		want := &target.GithubPush{
			Repo: &target.MockRepository{
				FullName: "duck8823/duci",
				URL:      "http://example.com",
			},
			Point: &github.SimpleTargetPoint{
				Ref: "test",
				SHA: random.String(16, random.Alphanumeric),
			},
		}
		want.SetGit(wantGit)

		// when
		got, err := target.NewGithubPush(want.Repo, want.Point)

		// then
		if err != nil {
			t.Errorf("must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got.Repo, want.Repo) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got.Repo, want.Repo))
		}

		if !cmp.Equal(got.GetGit(), want.GetGit()) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got.GetGit(), want.GetGit()))
		}
	})

	t.Run("when git not found", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := target.NewGithubPush(nil, nil)

		// then
		if err == nil {
			t.Errorf("must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})
}

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

		// and
		sut := &target.GithubPush{
			Repo:  repo,
			Point: point,
		}
		defer sut.SetGit(mockGit)()

		// when
		got, cleanup, err := sut.Prepare()
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

		// and
		sut := &target.GithubPush{
			Repo:  repo,
			Point: point,
		}
		defer sut.SetGit(mockGit)()

		// when
		got, cleanup, err := sut.Prepare()
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
