package logstore

import "github.com/duck8823/duci/infrastructure/store"

type StoreServiceImpl = storeServiceImpl

func (s *storeServiceImpl) SetDB(db store.Store) {
	s.db = db
}
