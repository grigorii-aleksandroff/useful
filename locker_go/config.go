package locker

import (
	"asiatix/internal/metrics"
	"time"
)

type Config struct {
	SlowlockTimeout time.Duration
	UnlockDelay     time.Duration
	KeyPrefix       string
}

type StoreConfig struct {
	*RedisConfig
	*MemoryConfig
}

type RedisConfig struct {
	Hosts   []string
	DbIndex int
}

type MemoryConfig struct{}

var lockerErrorMetric = metrics.NewUpCounterMetric("locker_error")
