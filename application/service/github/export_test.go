package github

type MockRepo struct {
	FullName string
	SSHURL   string
	CloneURL string
}

func (r *MockRepo) GetFullName() string {
	return r.FullName
}

func (r *MockRepo) GetSSHURL() string {
	return r.SSHURL
}

func (r *MockRepo) GetCloneURL() string {
	return r.CloneURL
}
