package job_service_test

import (
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/mock_job"
	"github.com/duck8823/duci/internal/container"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"os"
	"path"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	t.Run("with temporary directory", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// when
		err := job_service.Initialize(tmpDir)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("with invalid directory", func(t *testing.T) {
		// given
		tmpDir := path.Join("/path/to/invalid/dir")
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// when
		err := job_service.Initialize(tmpDir)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestGetInstance(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := job_service.GetInstance()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", err)
		}
	})

	t.Run("when instance is not nil", func(t *testing.T) {
		// given
		want := &job_service.StubService{
			ID: random.String(16, random.Alphanumeric),
		}

		// and
		container.Override(want)
		defer container.Clear()

		// when
		got, err := job_service.GetInstance()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})
}

func TestServiceImpl_FindBy(t *testing.T) {
	t.Run("when repo returns job", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		want := &job.Job{
			ID:       id,
			Finished: true,
			Stream: []job.LogLine{
				{
					Timestamp: time.Now(),
					Message:   "hello world",
				},
			},
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(want, nil)

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal but %+v", cmp.Diff(got, want))
		}
	})

	t.Run("when repo returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(nil, errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})
}

func TestServiceImpl_Start(t *testing.T) {
	t.Run("when repo returns nil", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: false})).
			Times(1).
			Return(nil)

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Start(id)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when repo returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: false})).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Start(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestServiceImpl_Append(t *testing.T) {
	t.Run("when failure find job", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())
		line := job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(nil, errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Append(id, line)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when failure save job", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())
		line := job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(&job.Job{ID: id, Finished: false}, nil)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: false, Stream: []job.LogLine{line}})).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Append(id, line)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when job not found", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())
		line := job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(nil, job.NotFound)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: false, Stream: []job.LogLine{line}})).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Append(id, line)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("without any error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())
		line := job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}

		// and
		stored := job.LogLine{Timestamp: time.Now(), Message: "Stored line"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(&job.Job{ID: id, Finished: false, Stream: []job.LogLine{stored}}, nil)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: false, Stream: []job.LogLine{stored, line}})).
			Times(1).
			Return(nil)

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Append(id, line)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})
}

func TestServiceImpl_Finish(t *testing.T) {
	t.Run("when find job, returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(nil, errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Finish(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when save, returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(&job.Job{ID: id, Finished: false}, nil)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: true})).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Finish(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("without any error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_job.NewMockRepository(ctrl)
		repo.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(&job.Job{ID: id, Finished: false}, nil)
		repo.EXPECT().
			Save(gomock.Eq(job.Job{ID: id, Finished: true})).
			Return(nil)

		// and
		sut := &job_service.ServiceImpl{}
		defer sut.SetRepo(repo)()

		// when
		err := sut.Finish(id)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})
}
