package initiator

import (
	"context"
	"log/slog"
	"time"

	"yadro.com/course/search/core"
)

type Initiator struct {
	log     *slog.Logger
	service *core.Service
	ttl     time.Duration
}

func NewInitiator(log *slog.Logger, service *core.Service, ttl time.Duration) *Initiator {
	return &Initiator{
		log:     log,
		service: service,
		ttl:     ttl,
	}
}

func (i *Initiator) Start(ctx context.Context) {
	i.log.Info("starting index initiator", "ttl", i.ttl)

	if err := i.service.BuildIndex(ctx); err != nil {
		i.log.Error("failed to build initial index", "error", err)
	}

	ticker := time.NewTicker(i.ttl)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				i.log.Info("stopping index initiator")
				return
			case <-ticker.C:
				i.log.Info("rebuilding index by ticker")
				if err := i.service.BuildIndex(ctx); err != nil {
					i.log.Error("failed to rebuild index", "error", err)
				}
			}
		}
	}()
}
