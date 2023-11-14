package background

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"time"
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

func (top TransactionalOutbox) Process(ctx context.Context) error {
	transactionalOutbox, err := top.transactionalOutboxRepository.FindByUnsent(ctx)
	if err != nil {
		return err
	}

	if transactionalOutbox.ID.NotExist() {
		return nil
	}

	event := domain.NewEvent(
		transactionalOutbox.Domain,
		transactionalOutbox.Type,
		string(transactionalOutbox.Body),
		time.Now().UTC(),
	)

	err = top.producer.Publish(ctx, event)
	if err != nil {
		return err
	}

	err = top.transactionalOutboxRepository.MarkToSent(ctx, transactionalOutbox.ID)
	if err != nil {
		return err
	}

	return nil
}
