package job

type Target interface {
	Prepare() (WorkDir, Cleanup, error)
}

type WorkDir string

func (w WorkDir) ToString() string {
	return string(w)
}

type Cleanup func() error
