package producer

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"log/slog"
)

type Broker interface {
	Publish(context.Context, any) error
}

type Producer struct {
	broker Broker
}

func NewProducer(b Broker) Producer {
	return Producer{
		broker: b,
	}
}

func (p Producer) Publish(ctx context.Context, event domain.Event) error {
	err := p.broker.Publish(ctx, event)
	if err != nil {
		logger.Slog.Info("Failed to publish event", slog.Any("error", err.Error()))
		return err
	}

	logger.Slog.Info("Event published successfully", slog.Any("event", event))

	return nil
}
