package job

type Target interface {
	Prepare() (WorkDir, Cleanup, error)
}

type WorkDir string

func (w WorkDir) String() string {
	return string(w)
}

type Cleanup func()
