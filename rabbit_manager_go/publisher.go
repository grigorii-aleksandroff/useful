package rabbit_manager_go

import (
	"asiatix/internal/crossdomain_events"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gobuffalo/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrEventNotFound = errors.New("event not found")

type Publisher struct {
	mu     sync.Mutex
	ch     *amqp.Channel
	conn   *amqp.Connection
	logger logger.FieldLogger
}

func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ch != nil {
		return p.ch.Close()
	}

	return nil
}

func (p *Publisher) PublishEvent(ctx context.Context, event crossdomain_events.CrossdomainEvent) error {
	if event.Retry() != nil {
		return p.publishEventWithRetry(ctx, event, *event.Retry())
	}
	return p.publish(ctx, event)
}

func (p *Publisher) publishEventWithRetry(ctx context.Context, event crossdomain_events.CrossdomainEvent, opts crossdomain_events.RetryOpts) (err error) {
	attempts := opts.Attempts
	delay := opts.Delay

	for i := 0; i < attempts; i++ {
		err = p.publish(ctx, event)
		if err == nil {
			return nil
		}

		p.logger.Warnf("publish attempt %d/%d failed for event %s: %v", i+1, attempts, event.Name().ToString(), err)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("all %d attempts failed for event %s: last error: %w", attempts, event.Name().ToString(), err)
}

func (p *Publisher) publish(ctx context.Context, event crossdomain_events.CrossdomainEvent) (err error) {
	p.mu.Lock()
	ch := p.ch
	p.mu.Unlock()

	if ch == nil {
		return fmt.Errorf("publisher channel is not initialized")
	}

	rule, ok := ruleset[event.Name()]
	if !ok {
		return fmt.Errorf("%w: %s", ErrEventNotFound, event.Name().ToString())
	}

	exchange := rule.Exchange.Name
	routingKey := rule.RoutingKey

	if rule.Queue != nil {
		if rule.Queue.Pattern != "" {
			dq, ok := event.(DynamicQueue)
			if !ok {
				return fmt.Errorf("event %s should implement DynamicQueue interface", event.Name().ToString())
			}

			queueName := dq.DynamicQueueName(rule.Queue.Pattern)
			routingKey = queueName

			if _, err = ch.QueueDeclare(
				queueName,
				rule.Queue.Durable,
				rule.Queue.AutoDelete,
				rule.Queue.Exclusive,
				rule.Queue.NoWait,
				rule.Queue.Args,
			); err != nil {
				return err
			}

			if err = ch.QueueBind(
				queueName,
				queueName,
				exchange,
				rule.Queue.NoWait,
				rule.Queue.Args,
			); err != nil {
				return err
			}
		} else if routingKey == "" {
			routingKey = rule.Queue.Name
		}
	}

	if routingKey == "" {
		return fmt.Errorf("routing key is required")
	}

	message, err := p.castMessage(event)
	if err != nil {
		return fmt.Errorf("cannot cast message from event: %w", err)
	}

	if err = ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		rule.Mandatory,
		rule.Immediate,
		*message,
	); err != nil {
		return fmt.Errorf("publish event %s error: %w", event.Name().ToString(), err)
	}

	p.logger.WithFields(map[string]interface{}{
		"exchange":     exchange,
		"routing_key":  routingKey,
		"content_type": message.ContentType,
		"body_size":    len(message.Body),
	}).Infof("published event %s to rabbit", event.Name().ToString())

	return nil
}

func (p *Publisher) castMessage(event crossdomain_events.CrossdomainEvent) (*amqp.Publishing, error) {
	payload, err := event.Payload()
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
	}

	switch v := payload.(type) {
	case *amqp.Publishing:
		return v, nil
	case []byte:
		msg.Body = v
	case string:
		msg.Body = []byte(v)
	default:
		return nil, fmt.Errorf("unsupported payload type %T", payload)
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	return &msg, nil
}
