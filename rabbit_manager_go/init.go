package rabbit_manager_go

import (
	"github.com/gobuffalo/logger"
	"sync"
)

var (
	manager *Manager
	once    sync.Once
)

func Init(c *Config, logger logger.FieldLogger) (*Manager, error) {
	var err error

	once.Do(func() {
		manager = newManager(c, logger)
		err = manager.connect()
		if err == nil {
			go manager.monitorConnection()
		}
	})

	return manager, err
}

func GetManager() *Manager {
	if manager == nil {
		panic("rabbit manager is not initialized")
	}
	return manager
}
