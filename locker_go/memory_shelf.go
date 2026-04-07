package locker

import (
	"time"
)

type memoryShelf struct {
	cells map[string]*cell
}

type cell struct {
	mutex          Mutex
	availableUntil time.Time
}

func newMemoryShelf() *memoryShelf {
	return &memoryShelf{
		cells: make(map[string]*cell),
	}
}

func (s *memoryShelf) getMutex(key string) (Mutex, bool) {
	mutexCell, exists := s.cells[key]
	if !exists {
		return nil, false
	}

	if !mutexCell.availableUntil.After(time.Now()) {
		s.deleteMutex(key)

		return nil, false
	}

	return mutexCell.mutex, true
}

func (s *memoryShelf) setMutex(key string, mutex Mutex, ttl time.Duration) {
	s.cells[key] = &cell{mutex, time.Now().Add(ttl).UTC()}
}

func (s *memoryShelf) deleteMutex(key string) {
	delete(s.cells, key)
}
