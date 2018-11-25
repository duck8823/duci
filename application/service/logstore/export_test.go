package logstore

import "github.com/duck8823/duci/infrastructure/store"

type StoreServiceImpl = storeServiceImpl

func (s *storeServiceImpl) SetDB(db store.Store) (reset func()) {
	tmp := s.db
	s.db = db
	return func() {
		s.db = tmp
	}
}
