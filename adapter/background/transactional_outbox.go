package background

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"log/slog"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
)

type TransactionalOutbox struct {
	producer                      domain.Producer
	transactionalOutboxRepository domain.TransactionalOutboxRepository
}

func NewTransactionalOutbox(
	producer domain.Producer,
	transactionalOutboxRepository domain.TransactionalOutboxRepository,
) TransactionalOutbox {
	return TransactionalOutbox{
		producer:                      producer,
		transactionalOutboxRepository: transactionalOutboxRepository,
	}
}

func (t TransactionalOutbox) Process(ctx context.Context) error {
	ctx, span := opentelemetry.NewSpan(ctx, "background.transactional_outbox.process")
	defer span.End()

	logger.Slog.Info("Starting processing transactional outbox")

	transactionalOutboxList, err := t.transactionalOutboxRepository.ListByUnsent(ctx)
	if err != nil {
		return err
	}

	if len(transactionalOutboxList) == 0 {
		logger.Slog.Info("No events found")
		logger.Slog.Info("Finishing processing transactional outbox")
		return nil
	}

	var (
		eventListToPublish   []domain.Event
		idListToMarkedToSent []valueobject.ID
	)

	for _, transactionalOutbox := range transactionalOutboxList {
		var event = domain.NewEvent(
			transactionalOutbox.Domain,
			transactionalOutbox.Type,
			string(transactionalOutbox.Body),
			time.Now().UTC(),
		)

		eventListToPublish = append(eventListToPublish, event)
		idListToMarkedToSent = append(idListToMarkedToSent, transactionalOutbox.ID)
	}

	err = t.producer.PublishBatch(ctx, eventListToPublish)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Info("Failed to publish event", slog.Any("error", err.Error()))
		return err
	}

	logger.Slog.Info("Events published", slog.Any("events", eventListToPublish))

	err = t.transactionalOutboxRepository.MarkToSentBatch(ctx, idListToMarkedToSent)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Info("Failed to mark event as sent", slog.Any("error", err.Error()))
		return err
	}

	logger.Slog.Info("Events marked as sent", slog.Any("ids", idListToMarkedToSent))

	logger.Slog.Info("Finishing processing transactional outbox")

	return nil
}
