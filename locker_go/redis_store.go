package locker

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"time"

	"asiatix/internal/redis_helpers"
)

type RedisStore struct {
	rs *redsync.Redsync
	*BaseStore
}

func NewRedisStore(config RedisConfig) *RedisStore {
	// CreateUniversalClient automatically detects cluster mode and handles DB configuration
	// - Multiple hosts: always cluster mode (DB not supported)
	// - Single host: checks if cluster, sets DB only for non-cluster
	client := redis_helpers.CreateUniversalClient(config.Hosts, config.DbIndex)

	if cmd := client.Ping(context.Background()); cmd.Err() != nil {
		panic(fmt.Sprintf("newRedisLocker error: %s", cmd.Err()))
	}

	pool := goredis.NewPool(client)

	return &RedisStore{
		rs:        redsync.New(pool),
		BaseStore: newBaseStore(),
	}
}

func (r *RedisStore) Type() StoreType {
	return Redis
}

func (r *RedisStore) NewMutex(key string, unlockDelay time.Duration) Mutex {
	return &RedsyncMutexAdapter{
		mutex: r.rs.NewMutex(key, redsync.WithExpiry(unlockDelay)),
	}
}

func (r *RedisStore) GetMutex(key string) (Mutex, bool) {
	return r.shelf.getMutex(key)
}

func (r *RedisStore) SetMutex(key string, mutex Mutex, unlockDelay time.Duration) {
	r.shelf.setMutex(key, mutex, unlockDelay)
}

func (r *RedisStore) DeleteMutex(key string) {
	r.shelf.deleteMutex(key)
}

type RedsyncMutexAdapter struct {
	mutex *redsync.Mutex
}

func (r *RedsyncMutexAdapter) Lock() error {
	return r.mutex.Lock()
}

func (r *RedsyncMutexAdapter) Unlock() (bool, error) {
	return r.mutex.Unlock()
}
