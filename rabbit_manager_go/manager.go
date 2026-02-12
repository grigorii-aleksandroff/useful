package rabbit_manager_go

import (
	"context"
	"fmt"
	"github.com/gobuffalo/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"net"
	"sync"
	"time"
)

type relistenable interface {
	relisten()
}

type Manager struct {
	mu     sync.Mutex
	conf   *Config
	url    string
	ctx    context.Context
	conn   *amqp.Connection
	cancel context.CancelFunc
	logger logger.FieldLogger

	p         *Publisher
	consumers []relistenable

	reconnecting chan struct{}
}

func newManager(c *Config, logger logger.FieldLogger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		conf:         c,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
		reconnecting: make(chan struct{}, 1),
	}
	return m
}

func (m *Manager) Publisher() *Publisher {
	return m.p
}

func (m *Manager) monitorConnection() {
	ticker := time.NewTicker(m.conf.ReconnectDelay)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			conn := m.conn
			m.mu.Unlock()
			if conn == nil || conn.IsClosed() {
				m.tryReconnect()
			}
		}
	}
}

func (m *Manager) connect() error {
	conn, err := connectRabbit(m.connectionUrl(), m.conf.ConnectionAttempt, m.conf.ConnectionTimeout, m.logger)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	defer ch.Close()

	if err := declareRules(ch); err != nil {
		conn.Close()
		return fmt.Errorf("failed to declare topology: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.conn != nil && !m.conn.IsClosed() {
		_ = m.conn.Close()
	}

	pubCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	m.conn = conn
	if m.p == nil {
		m.p = &Publisher{
			conn:   conn,
			ch:     pubCh,
			logger: m.logger,
		}
	} else {
		m.p.mu.Lock()
		if m.p.ch != nil {
			_ = m.p.ch.Close()
		}
		m.p.conn = conn
		m.p.ch = pubCh
		m.p.mu.Unlock()
	}

	return nil
}

func (m *Manager) tryReconnect() {
	select {
	case m.reconnecting <- struct{}{}:
		go func() {
			defer func() {
				<-m.reconnecting
			}()
			m.reconnect()
		}()
	default:
	}
}

func (m *Manager) reconnect() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		if err := m.connect(); err != nil {
			m.logger.Warnf("rabbit reconnect failed: %v, retry in %v", err, m.conf.ReconnectDelay)

			select {
			case <-m.ctx.Done():
				return
			case <-time.After(m.conf.ReconnectDelay):
			}

			continue
		}

		m.mu.Lock()
		consumers := append([]relistenable(nil), m.consumers...)
		m.mu.Unlock()

		for _, c := range consumers {
			c.relisten()
		}

		m.logger.Infof("rabbit: reconnected successfully")
		return
	}
}

func (m *Manager) connectionUrl() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s/%s",
		m.conf.User,
		m.conf.Password,
		net.JoinHostPort(m.conf.Host, m.conf.Port),
		m.conf.Vhost,
	)
}

func connectRabbit(url string, attempt int, timeout time.Duration, logger logger.Logger) (*amqp.Connection, error) {
	var lastErr error

	for i := 0; i < attempt; i++ {
		conn, err := amqp.Dial(url)
		if err != nil {
			lastErr = err
			logger.Warnf("rabbit: failed to connect (attempt %d): %v", i+1, err)
			time.Sleep(timeout)
			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("rabbit: all connection attempts failed: %w", lastErr)
}

func declareRules(ch *amqp.Channel) error {
	for _, rule := range ruleset {
		if rule.Exchange == nil {
			return fmt.Errorf("rule exchange is nil: exchange required")
		}

		if err := ch.ExchangeDeclare(
			rule.Exchange.Name,
			rule.Exchange.Kind.ToString(),
			rule.Exchange.Durable,
			rule.Exchange.AutoDelete,
			rule.Exchange.Internal,
			rule.Exchange.NoWait,
			rule.Exchange.Args,
		); err != nil {
			return fmt.Errorf("declare exchange %s: %w", rule.Exchange.Name, err)
		}

		if rule.Queue == nil {
			return nil
		}

		if rule.Queue.Pattern == "" {
			q, err := ch.QueueDeclare(
				rule.Queue.Name,
				rule.Queue.Durable,
				rule.Queue.AutoDelete,
				rule.Queue.Exclusive,
				rule.Queue.NoWait,
				rule.Queue.Args,
			)
			if err != nil {
				return fmt.Errorf("declare queue %s: %w", rule.Queue.Name, err)
			}

			if err := ch.QueueBind(
				q.Name,
				rule.Queue.BindKey,
				rule.Exchange.Name,
				rule.Queue.NoWait,
				rule.Queue.Args,
			); err != nil {
				return fmt.Errorf("bind queue %s: %w", q.Name, err)
			}
		}
	}

	return nil
}
