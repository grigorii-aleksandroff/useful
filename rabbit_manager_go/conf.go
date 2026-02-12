package rabbit_manager_go

import "time"

type Config struct {
	Host              string
	Port              string
	User              string
	Password          string
	Vhost             string
	ConnectionAttempt int
	ConnectionTimeout time.Duration
	ReconnectDelay    time.Duration
}
