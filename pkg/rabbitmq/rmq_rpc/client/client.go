package client

import (
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

var ErrConnectionClosed = errors.New("rmq_rpc client - Client - RemoteCall - Connection closed")

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second

	Success = "success"
)

// Структура, представляющая сообщение RabbitMQ
type Message struct {
	Queue         string
	Priority      uint8
	ContentType   string
	Body          []byte
	ReplyTo       string
	CorrelationID string
}

// Структура, хранящая информацию о вызове, который ожидает ответа
type pendingCall struct {
	done   chan struct{}
	status string
	body   []byte
}

// RPC-клиента
type Client struct {
	conn           *rmq_rpc.Connection
	serverExchange string
	error          chan error
	stop           chan struct{}

	rw    sync.RWMutex
	calls map[string]*pendingCall

	timeout time.Duration
}

// Создание клиента
func NewRabbitMQClient(url, serverExchange, clientExchange string, opts ...Option) (*Client, error) {
	cfg := rmq_rpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		conn:           rmq_rpc.NewConnectionRabbitMQ(clientExchange, cfg),
		serverExchange: serverExchange,
		error:          make(chan error),
		stop:           make(chan struct{}),
		calls:          make(map[string]*pendingCall),
		timeout:        _defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc client - NewClient - c.conn.AttemptConnect: %w", err)
	}

	// Прослушивание сообщений
	go c.consumer()

	return c, nil
}

// Публикует сообщение в RabbitMQ
func (c *Client) publish(corrID, handler string, request interface{}) error {
	var (
		requestBody []byte
		err         error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("publish - json.Marshal: %w", err)
		}
	}

	err = c.conn.Channel.Publish(c.serverExchange, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.ConsumerExchange,
			Type:          handler,
			Body:          requestBody,
		})
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (c *Client) RemoteCall(ctx context.Context, handler string, request, response interface{}) error {
	// Если клиент не завершился
	select {
	case <-c.stop:
		time.Sleep(c.timeout)
		select {
		case <-c.stop:
			return ErrConnectionClosed
		default:
		}
	default:
	}

	corrID := uuid.New().String()

	err := c.publish(corrID, handler, request)
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.publish: %w", err)
	}

	call := &pendingCall{done: make(chan struct{})}

	// Добавляет вызов в c.calls, чтобы отслеживать его статус.
	c.addCall(corrID, call)
	defer c.deleteCall(corrID)

	// Если не выполняется за определенное время то ошибка
	select {
	case <-time.After(c.timeout):
		return rmq_rpc.ErrTimeout
	case <-ctx.Done():
		return ctx.Err() //nolint:wrapcheck // just sending raw error
	case <-call.done:
	}

	return c.handleCallStatus(call, response)
}

func (c *Client) handleCallStatus(call *pendingCall, response interface{}) error {
	if call.status == Success {
		err := json.Unmarshal(call.body, &response)
		if err != nil {
			return fmt.Errorf("rmq_rpc client - Client - handleCallStatus - json.Unmarshal: %w", err)
		}

		return nil
	}

	if call.status == rmq_rpc.ErrBadHandler.Error() {
		return rmq_rpc.ErrBadHandler
	}

	if call.status == rmq_rpc.ErrNotFound.Error() {
		return rmq_rpc.ErrNotFound
	}

	return fmt.Errorf("%w: %v", rmq_rpc.ErrCallStatus, call.status)
}

func (c *Client) consumer() {
	for {
		select {
		case <-c.stop:
			return
		case d, opened := <-c.conn.Delivery:
			if !opened {
				c.reconnect()

				return
			}

			_ = d.Ack(false)

			c.getCall(&d)
		}
	}
}

func (c *Client) reconnect() {
	close(c.stop)

	err := c.conn.AttemptConnect()
	if err != nil {
		c.error <- err
		close(c.error)

		return
	}

	c.stop = make(chan struct{})

	go c.consumer()
}

func (c *Client) getCall(d *amqp.Delivery) {
	c.rw.RLock()
	call, ok := c.calls[d.CorrelationId]
	c.rw.RUnlock()

	if !ok {
		return
	}

	call.status = d.Type
	call.body = d.Body
	close(call.done)
}

func (c *Client) addCall(corrID string, call *pendingCall) {
	c.rw.Lock()
	c.calls[corrID] = call
	c.rw.Unlock()
}

func (c *Client) deleteCall(corrID string) {
	c.rw.Lock()
	delete(c.calls, corrID)
	c.rw.Unlock()
}

func (c *Client) Notify() <-chan error {
	return c.error
}

func (c *Client) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - Shutdown - c.Connection.Close: %w", err)
	}

	return nil
}
