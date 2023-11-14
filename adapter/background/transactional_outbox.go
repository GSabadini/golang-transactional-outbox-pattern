package background

import (
	"context"
	"encoding/json"
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

	event, err := top.buildEvent(transactionalOutbox)
	if err != nil {
		return err
	}

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

func (top TransactionalOutbox) buildEvent(transactionalOutbox domain.TransactionalOutbox) (domain.Event, error) {
	var transaction domain.Transaction
	err := json.Unmarshal(transactionalOutbox.Body, &transaction)
	if err != nil {
		return domain.Event{}, err
	}

	return domain.NewEvent(
		transactionalOutbox.Domain,
		transactionalOutbox.Type,
		string(transactionalOutbox.Body),
		time.Now().UTC(),
	), nil
}
