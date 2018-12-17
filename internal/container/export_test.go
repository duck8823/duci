package container

type SingletonContainer = singletonContainer

func (s *SingletonContainer) SetValues(values map[string]interface{}) (reset func()) {
	tmp := s.values
	s.values = values
	return func() {
		s.values = tmp
	}
}

func (s *SingletonContainer) GetValues() map[string]interface{} {
	return s.values
}

func SetInstance(ins *singletonContainer) (reset func()) {
	tmp := instance
	instance = ins
	return func() {
		instance = tmp
	}
}
