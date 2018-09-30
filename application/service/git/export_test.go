package git

import "gopkg.in/src-d/go-git.v4"

func SetPlainCloneFunc(f func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)) {
	plainClone = f
}
