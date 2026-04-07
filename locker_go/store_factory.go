package locker

import "fmt"

func GetStore(storeType StoreType, config StoreConfig) (Store, error) {
	switch storeType {
	case Redis:
		if config.RedisConfig == nil {
			return nil, fmt.Errorf("config for redis should not be nil")
		}

		return NewRedisStore(*config.RedisConfig), nil
	case Memory:
		if config.MemoryConfig == nil {
			return nil, fmt.Errorf("config for redis should not be nil")
		}

		return NewMemoryStore(*config.MemoryConfig), nil
	default:
		return NewMemoryStore(*config.MemoryConfig), nil
	}
}
