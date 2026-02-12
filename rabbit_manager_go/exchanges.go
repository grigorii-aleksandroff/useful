package rabbit_manager_go

import amqp "github.com/rabbitmq/amqp091-go"

const (
	DefaultExchangeName = "default_exchange"
)

type ExchangeKind string

func (e ExchangeKind) ToString() string {
	return string(e)
}

const (
	ExchangeKindDirect  ExchangeKind = `direct`
	ExchangeKindTopic   ExchangeKind = `topic`
	ExchangeKindFanout  ExchangeKind = `fanout`
	ExchangeKindHeaders ExchangeKind = `headers`
)

type Exchange struct {
	Name       string
	Kind       ExchangeKind
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

var DefaultExchange = &Exchange{
	Name:    DefaultExchangeName,
	Kind:    ExchangeKindFanout,
	Durable: true,
}
