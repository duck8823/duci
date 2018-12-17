package target

type MockRepository struct {
	FullName string
	URL      string
}

func (r *MockRepository) GetFullName() string {
	return r.FullName
}

func (r *MockRepository) GetSSHURL() string {
	return r.URL
}

func (r *MockRepository) GetCloneURL() string {
	return r.URL
}
