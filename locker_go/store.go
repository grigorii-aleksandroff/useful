package locker

import (
	"sync"
	"time"
)

type BaseStore struct {
	shelf *memoryShelf
	*sync.Mutex
}

func newBaseStore() *BaseStore {
	return &BaseStore{
		shelf: newMemoryShelf(),
		Mutex: new(sync.Mutex),
	}
}

type Store interface {
	Type() StoreType
	NewMutex(key string, unlockDelay time.Duration) Mutex
	GetMutex(key string) (Mutex, bool)
	SetMutex(key string, mutex Mutex, unlockDelay time.Duration)
	DeleteMutex(key string)
	Lock()
	Unlock()
}

type Mutex interface {
	Lock() error
	Unlock() (bool, error)
}
