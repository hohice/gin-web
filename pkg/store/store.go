package store

import "sync"

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	Close()
}

var (
	storeMethod Store
)

var locker = new(sync.Mutex)

func SetMethodInstance(store Store) {
	locker.Lock()
	defer locker.Unlock()
	storeMethod = store
}

func GetMethodInstance() Store {
	locker.Lock()
	defer locker.Unlock()
	return storeMethod
}
