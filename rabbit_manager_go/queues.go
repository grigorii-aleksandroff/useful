package rabbit_manager_go

import amqp "github.com/rabbitmq/amqp091-go"

const (
	DefaultQueueName = "default"
)

const (
	QueueArgMessageTTL           = `x-message-ttl`
	QueueArgExpires              = `x-expires`
	QueueArgMaxLength            = `x-max-length`
	QueueArgMaxLengthBytes       = `x-max-length-bytes`
	QueueArgOverflow             = `x-overflow`
	QueueArgDeadLetterExchange   = `x-dead-letter-exchange`
	QueueArgDeadLetterRoutingKey = `x-dead-letter-routing-key`
	QueueArgMaxPriority          = `x-max-priority`
	QueueArgQueueMode            = `x-queue-mode`
	QueueArgQueueMaster          = `x-queue-master`
	QueueArgHaPolicy             = `x-ha-policy`
	QueueArgHaNodes              = `x-ha-nodes`
)

type DynamicQueue interface {
	DynamicQueueName(pattern string) string
}

type Queue struct {
	Name       string // uses for static queue
	Pattern    string // uses for dynamic queues (operation.*.created, operation.*.updated etc.)
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	BindKey    string
}

var DefaultQueue = &Queue{
	Name:    DefaultQueueName,
	Durable: true,
}
