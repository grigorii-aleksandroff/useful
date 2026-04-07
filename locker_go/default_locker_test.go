package locker

import (
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestKeyLock struct {
	Key string
}

func (k TestKeyLock) GetKeyLockString() (string, error) {
	return "test-key", nil
}

func TestLocker_shelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	unlockDelay := 100 * time.Millisecond
	shelf := newMemoryShelf()
	mutex := NewMockMutex(ctrl)

	shelf.setMutex("test-period", mutex, unlockDelay)

	_, exist := shelf.getMutex("test-period")
	assert.Equal(t, true, exist)

	time.Sleep(unlockDelay / 2)
	_, exist = shelf.getMutex("test-period")
	assert.Equal(t, true, exist)

	time.Sleep(unlockDelay / 2)
	_, exist = shelf.getMutex("test-period")
	assert.Equal(t, false, exist)
}

// TODO:: реализовать мок для редиса и передать его в redsync.New(redis) для того чтобы гибко проверять работу redsync
func TestLocker_mock_deadlock(t *testing.T) {
	var callCount int
	unlockDelay := 20 * time.Second
	slowlockTimeout := 5 * time.Second

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	mockMutex := NewMockMutex(ctrl)

	mockMutex.EXPECT().Lock().Return(nil).AnyTimes()
	mockMutex.EXPECT().Unlock().Return(true, nil).AnyTimes()

	mockStore.EXPECT().Type().Return(Mock).AnyTimes()
	mockStore.EXPECT().NewMutex(gomock.Any(), gomock.Any()).Return(mockMutex).AnyTimes()
	mockStore.EXPECT().GetMutex(gomock.Any()).DoAndReturn(func(key string) (Mutex, bool) {
		callCount++
		if callCount%5 == 1 {
			return nil, false
		}
		return mockMutex, true
	}).AnyTimes()

	mockStore.EXPECT().SetMutex(gomock.Any(), gomock.Any(), unlockDelay).AnyTimes()
	mockStore.EXPECT().DeleteMutex(gomock.Any()).AnyTimes()
	mockStore.EXPECT().Lock().Return().AnyTimes()
	mockStore.EXPECT().Unlock().Return().AnyTimes()

	lk := newDefaultLocker(Config{
		SlowlockTimeout: slowlockTimeout,
		UnlockDelay:     unlockDelay,
		KeyPrefix:       "lkr-test",
	}, mockStore)

	assert.Equal(t, Mock, mockStore.Type(), "type should be mock")
	startDeadLockAndConcurrencyTests(t, lk)
}

func TestLocker_memory_deadlock(t *testing.T) {
	unlockDelay := 20 * time.Second
	slowlockTimeout := 5 * time.Second
	store := NewMemoryStore(MemoryConfig{})

	lk := newDefaultLocker(Config{
		SlowlockTimeout: slowlockTimeout,
		UnlockDelay:     unlockDelay,
		KeyPrefix:       "lkr-test",
	}, store)

	assert.Equal(t, Memory, store.Type(), "type should be memory")
	startDeadLockAndConcurrencyTests(t, lk)
}

func startDeadLockAndConcurrencyTests(t *testing.T, lk Locker) {
	testKeys := []string{"key1", "key2", "key3"}

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := TestKeyLock{Key: testKeys[idx%len(testKeys)]}
			err := lk.Lock(key)
			assert.NoError(t, err, "gorutine shouldn't get error after the unlock")
			time.Sleep(100 * time.Millisecond)
			err = lk.Unlock(key)
			assert.NoError(t, err, "gorutine should unlock without blocks")
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		key := TestKeyLock{Key: "key1"}
		err := lk.Lock(key)
		assert.NoError(t, err, "locking key1 should complete without errors")
		time.Sleep(1000 * time.Millisecond)
		err = lk.Unlock(key)
		assert.NoError(t, err, "unlocking key1 should complete without errors")
	}()

	wg.Wait()
}
