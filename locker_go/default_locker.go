package locker

import (
	"errors"
	"fmt"
	"time"
)

var _ Locker = (*DefaultLocker)(nil)

type DefaultLocker struct {
	store  Store
	config Config
}

func newDefaultLocker(lockerConfig Config, store Store) Locker {
	return &DefaultLocker{
		config: lockerConfig,
		store:  store,
	}
}

func (l *DefaultLocker) Lock(keyLock KeyLock) error {
	var errString = "default locker lock error: %s, %w"

	preparedKey, err := prepareKey(l.config.KeyPrefix, keyLock)
	if err != nil {
		return fmt.Errorf(errString, "empty", err)
	}

	l.store.Lock()

	if _, ok := l.store.GetMutex(preparedKey); ok {
		l.store.Unlock()

		if err = l.slowLock(preparedKey); err != nil {
			return fmt.Errorf(errString, preparedKey, err)
		}

		return nil
	}

	mutex := l.store.NewMutex(preparedKey, l.config.UnlockDelay)
	if err = mutex.Lock(); err != nil {
		l.store.Unlock()

		return fmt.Errorf(errString, preparedKey, err)
	}

	l.store.SetMutex(preparedKey, mutex, l.config.UnlockDelay)
	l.store.Unlock()

	return nil
}

func (l *DefaultLocker) Unlock(keyLock KeyLock) error {
	var errString = "default locker unlock error: %s, %w"

	preparedKey, err := prepareKey(l.config.KeyPrefix, keyLock)
	if err != nil {
		return fmt.Errorf(errString, "empty", err)
	}

	l.store.Lock()
	defer l.store.Unlock()

	mutex, ok := l.store.GetMutex(preparedKey)
	if !ok {
		return fmt.Errorf(errString, preparedKey, errors.New("key not found"))
	}
	if mutex == nil {
		return fmt.Errorf(errString, preparedKey, errors.New("mutex is nil"))
	}

	ok, err = mutex.Unlock()
	if err != nil || !ok {
		return fmt.Errorf(errString, preparedKey, err)
	}

	l.store.DeleteMutex(preparedKey)

	return nil
}

func (l *DefaultLocker) slowLock(preparedKey string) error {
	timeout := time.After(l.config.SlowlockTimeout)
	ticker := time.NewTicker(10 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("slow lock timeout for key: %s", preparedKey)
		case <-ticker.C:
			l.store.Lock()

			if _, exists := l.store.GetMutex(preparedKey); exists {
				l.store.Unlock()
				continue
			}

			mutex := l.store.NewMutex(preparedKey, l.config.UnlockDelay)

			if err := mutex.Lock(); err != nil {
				l.store.Unlock()
				return fmt.Errorf("slow lock error: %s, %w", preparedKey, err)
			}

			l.store.SetMutex(preparedKey, mutex, l.config.UnlockDelay)
			l.store.Unlock()

			return nil
		}
	}
}
