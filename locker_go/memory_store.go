package locker

import (
	"sync"
	"time"
)

type MemoryStore struct {
	*BaseStore
}

func NewMemoryStore(config MemoryConfig) *MemoryStore {
	return &MemoryStore{
		BaseStore: newBaseStore(),
	}
}

func (m *MemoryStore) Type() StoreType {
	return Memory
}

func (m *MemoryStore) NewMutex(key string, unlockDelay time.Duration) Mutex {
	return &SyncMutexAdapter{mutex: &sync.Mutex{}}
}

func (m *MemoryStore) GetMutex(key string) (Mutex, bool) {
	return m.shelf.getMutex(key)
}

func (m *MemoryStore) SetMutex(key string, mutex Mutex, unlockDelay time.Duration) {
	m.shelf.setMutex(key, mutex, unlockDelay)
}

func (m *MemoryStore) DeleteMutex(key string) {
	m.shelf.deleteMutex(key)
}

type SyncMutexAdapter struct {
	mutex *sync.Mutex
}

func (s *SyncMutexAdapter) Lock() error {
	s.mutex.Lock()
	return nil
}

func (s *SyncMutexAdapter) Unlock() (bool, error) {
	s.mutex.Unlock()
	return true, nil
}
