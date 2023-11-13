package producer

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"log/slog"

	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
)

type BrokerPublish interface {
	Publish()
}

type Publish struct {
	BrokerPublish BrokerPublish
}

func NewPublish() Publish {
	return Publish{}
}

func (p Publish) Publish(ctx context.Context, event domain.Event) error {
	logger.Slog.Info("Publish event", slog.Any("event", event))

	p.BrokerPublish.Publish()

	return nil
}
