package eventbus

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"yadro.com/course/search/core"
)

type Subscriber struct {
	nc      *nats.Conn
	log     *slog.Logger
	service *core.Service
}

func NewSubscriber(brokerAddress string, log *slog.Logger, service *core.Service) (*Subscriber, error) {
	if brokerAddress == "" {
		return nil, fmt.Errorf("broker address is empty")
	}
	nc, err := nats.Connect(brokerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}
	return &Subscriber{
		nc:      nc,
		log:     log,
		service: service,
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context) error {
	_, err := s.nc.Subscribe("xkcd.db.updated", func(msg *nats.Msg) {
		s.log.Info("received update event, rebuilding index")

		if err := s.service.BuildIndex(context.Background()); err != nil {
			s.log.Error("failed to rebuild index", "error", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to updates: %w", err)
	}
	s.log.Info("subscribed to xkcd.db.updated")

	if err := s.nc.Flush(); err != nil {
		return fmt.Errorf("failed to flush nats connection: %w", err)
	}
	return nil
}

func (s *Subscriber) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
