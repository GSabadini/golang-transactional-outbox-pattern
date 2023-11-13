package background

import (
	"context"
	"encoding/json"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
)

type TransactionalOutbox struct {
	producer                   domain.Producer
	transactionalOutboxFinder  domain.TransactionalOutboxFinder
	transactionalOutboxUpdater domain.TransactionalOutboxUpdater
}

func NewTransactionalOutbox(
	producer domain.Producer,
	transactionalOutboxFinder domain.TransactionalOutboxFinder,
	transactionalOutboxUpdater domain.TransactionalOutboxUpdater,
) TransactionalOutbox {
	return TransactionalOutbox{
		producer:                   producer,
		transactionalOutboxFinder:  transactionalOutboxFinder,
		transactionalOutboxUpdater: transactionalOutboxUpdater,
	}
}

func (top TransactionalOutbox) Process(ctx context.Context) error {
	transactionOutbox, err := top.transactionalOutboxFinder.FindByUnsent(ctx)
	if err != nil {
		return err
	}

	if transactionOutbox.ID.NotExist() {
		return nil
	}

	var transaction domain.Transaction
	err = json.Unmarshal(transactionOutbox.Body, &transaction)
	if err != nil {
		return err
	}

	err = top.producer.Publish(ctx, domain.Event{})
	if err != nil {
		return err
	}

	err = top.transactionalOutboxUpdater.MarkToSent(ctx, transactionOutbox.ID)
	if err != nil {
		return err
	}

	return nil
}
