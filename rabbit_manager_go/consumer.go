package rabbit_manager_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sort"
	"sync"
)

type ConsumerHandler[T any] func(
	ctx context.Context,
	log logger.FieldLogger,
	data T,
) error

type Subscriber[T any] interface {
	Options() ConsumerOptions
	Queue() Queue
	CastData(msg json.RawMessage) (T, error)
}

type ConsumerOptions struct {
	Tag       string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type handlerEntry[T any] struct {
	Priority int
	Handler  ConsumerHandler[T]
}

type Consumer[T any] struct {
	ctx        context.Context
	getConn    func() *amqp.Connection
	handlers   []handlerEntry[T]
	logger     logger.FieldLogger
	subscriber Subscriber[T]

	mu      sync.Mutex
	cancel  context.CancelFunc
	running bool

	errCh chan error
}

func NewConsumer[T any](subscriber Subscriber[T]) (*Consumer[T], error) {
	if manager == nil {
		return nil, fmt.Errorf("rabbit manager is not initialized")
	}

	c := &Consumer[T]{
		ctx:        manager.ctx,
		logger:     manager.logger,
		subscriber: subscriber,
		errCh:      make(chan error, 10),
		getConn: func() *amqp.Connection {
			manager.mu.Lock()
			defer manager.mu.Unlock()
			return manager.conn
		},
	}

	manager.mu.Lock()
	manager.consumers = append(manager.consumers, c)
	manager.mu.Unlock()

	return c, nil
}

func (c *Consumer[T]) Define(handler ConsumerHandler[T], priority int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handlerEntry[T]{
		Priority: priority,
		Handler:  handler,
	})

	sort.SliceStable(c.handlers, func(i, j int) bool {
		return c.handlers[i].Priority > c.handlers[j].Priority
	})

	c.logger.Debugf(
		"consumer: new handler for queue %s with priority %d",
		c.subscriber.Queue().Name,
		priority,
	)
}

func (c *Consumer[T]) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}

	c.running = true
	localCtx, cancel := context.WithCancel(c.ctx)
	c.cancel = cancel

	if c.errCh == nil {
		c.errCh = make(chan error, 10)
	}

	go c.listen(localCtx)
}

func (c *Consumer[T]) Stop() {
	c.mu.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	c.running = false
	c.mu.Unlock()
}

func (c *Consumer[T]) listen(ctx context.Context) {
	ch, err := c.getConn().Channel()
	if err != nil {
		select {
		case c.errCh <- err:
		default:
		}
		return
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		c.subscriber.Queue().Name,
		c.subscriber.Queue().Durable,
		c.subscriber.Queue().AutoDelete,
		c.subscriber.Queue().Exclusive,
		c.subscriber.Queue().NoWait,
		c.subscriber.Queue().Args,
	)
	if err != nil {
		c.errCh <- err
		return
	}

	msgs, err := ch.ConsumeWithContext(
		ctx,
		c.subscriber.Queue().Name,
		c.subscriber.Options().Tag,
		c.subscriber.Options().AutoAck,
		c.subscriber.Options().Exclusive,
		c.subscriber.Options().NoLocal,
		c.subscriber.Options().NoWait,
		c.subscriber.Options().Args,
	)
	if err != nil {
		select {
		case c.errCh <- err:
		default:
		}
		return
	}

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				c.logger.Warnf("consumer: delivery channel closed for queue %s", c.subscriber.Queue().Name)
				return
			}

			data, err := c.subscriber.CastData(d.Body)
			if err != nil {
				c.logger.Errorf("error decoding message from %s: %v", c.subscriber.Queue().Name, err)
				if !c.subscriber.Options().AutoAck {
					_ = d.Nack(false, true)
				}
				continue
			}

			for _, h := range c.handlers {
				if err := h.Handler(ctx, c.logger, data); err != nil {
					c.logger.Errorf("handler error for queue %s: %v", c.subscriber.Queue().Name, err)
				}
			}

			if !c.subscriber.Options().AutoAck {
				_ = d.Ack(false)
			}

		case <-ctx.Done():
			c.logger.Infof("consumer: context canceled for queue %s", c.subscriber.Queue().Name)
			return
		}
	}
}

func (c *Consumer[T]) relisten() {
	c.Stop()
	c.Start()
}
