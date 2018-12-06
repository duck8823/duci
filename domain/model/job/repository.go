package job

// Repository is Job Repository
type Repository interface {
	Get(ID) (*Job, error)
	Start(ID) error
	Append(ID, LogLine) error
	Finish(ID) error
}
