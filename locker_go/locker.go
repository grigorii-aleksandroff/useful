package locker

import (
	"asiatix/internal/metrics"
	"fmt"
	"sync"
)

var (
	AppLocker Locker
	mux       sync.Mutex
	once      sync.Once
)

type Locker interface {
	Lock(keyLock KeyLock) error
	Unlock(keyLock KeyLock) error
}

type KeyLock interface {
	GetKeyLockString() (string, error)
}

func GetAppLocker() (Locker, error) {
	if AppLocker == nil {
		return nil, fmt.Errorf("no application locker initialized")
	}
	return AppLocker, nil
}

func Init(config Config, store Store, exporter metrics.Exporter) {
	once.Do(func() {
		if AppLocker == nil {
			mux.Lock()
			defer mux.Unlock()

			// добавить фабрику после добавления других реализаций локера
			locker := newDefaultLocker(config, store)

			AppLocker = newLockerWrapper(locker, exporter)
		}
	})
}
