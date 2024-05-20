package server

import (
	"WalletRieltaTestTask/pkg/logger"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log/slog"
	"time"
)

const (
	_defaultWaitTime        = 2 * time.Second
	_defaultAttempts        = 10
	_defaultTimeout         = 2 * time.Second
	_defaultGoroutinesCount = 24

	Success = "success"
)

type CallHandler func(*amqp.Delivery) (interface{}, error)

type Server struct {
	conn   *rmq_rpc.Connection
	error  chan error
	stop   chan struct{}
	router map[string]CallHandler

	timeout         time.Duration
	goroutinesCount int

	logger *slog.Logger
}

func New(url, serverExchange string, router map[string]CallHandler, l *slog.Logger, opts ...Option) (*Server, error) {
	cfg := rmq_rpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		conn:            rmq_rpc.NewConnectionRabbitMQ(serverExchange, cfg),
		error:           make(chan error),
		stop:            make(chan struct{}),
		router:          router,
		timeout:         _defaultTimeout,
		goroutinesCount: _defaultGoroutinesCount,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err)
	}

	return s, nil
}

func (s *Server) MustRun() {
	for i := 0; i < s.goroutinesCount; i++ {
		go s.consumer()
	}
}

func (s *Server) consumer() {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-s.conn.Delivery:
			if !opened {
				s.reconnect()

				return
			}

			_ = d.Ack(false)

			s.serveCall(&d)
		}
	}
}

func (s *Server) serveCall(d *amqp.Delivery) {
	callHandler, ok := s.router[d.Type]
	if !ok {
		s.publish(d, nil, rmq_rpc.ErrBadHandler.Error())

		return
	}

	response, err := callHandler(d)
	if err != nil {
		s.publish(d, nil, err.Error())

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		s.logger.Error("rmq_rpc server - Server - serveCall - json.Marshal", logger.Err(err))
	}

	s.publish(d, body, Success)
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.Channel.Publish(d.ReplyTo, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		})
	if err != nil {
		s.logger.Error("rmq_rpc server - Server - publish - s.conn.Channel.Publish", logger.Err(err))
	}
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect()
	if err != nil {
		s.error <- err
		close(s.error)

		return
	}

	s.stop = make(chan struct{})

	go s.consumer()
}

func (s *Server) Notify() <-chan error {
	return s.error
}

func (s *Server) Shutdown() error {
	select {
	case <-s.error:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.timeout)

	err := s.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc server - Server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}
