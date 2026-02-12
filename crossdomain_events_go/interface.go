package crossdomain_events

import (
	"context"
	"time"
)

type Broker string

type EventName string

func (e EventName) ToString() string {
	return string(e)
}

type RetryOpts struct {
	Attempts int
	Delay    time.Duration
}

type Publisher interface {
	PublishEvent(ctx context.Context, event CrossdomainEvent) error
}

type CrossdomainEvent interface {
	Name() EventName
	Broker() Broker
	Retry() *RetryOpts
	Payload() (interface{}, error)
}
