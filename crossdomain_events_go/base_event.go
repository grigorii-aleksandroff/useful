package crossdomain_events

type BaseEvent struct {
	retryOpts *RetryOpts
}

func (e *BaseEvent) Retry() *RetryOpts {
	return e.retryOpts
}

func (e *BaseEvent) WithRetry(opts RetryOpts) {
	e.retryOpts = &opts
}

func (e *BaseEvent) Name() EventName {
	return DefaultEventName
}

func (e *BaseEvent) Broker() Broker {
	return RabbitBroker
}

func (e *BaseEvent) Payload() (interface{}, error) {
	return nil, nil
}
