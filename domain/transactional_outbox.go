package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
)

type (
	TransactionalOutboxRepository interface {
		Create(context.Context, *sql.Tx, TransactionalOutbox) error
		ListByUnsent(context.Context) ([]TransactionalOutbox, error)
		MarkToSent(context.Context, valueobject.ID) error
		MarkToSentBatch(context.Context, []valueobject.ID) error
	}

	TransactionalOutbox struct {
		ID        valueobject.ID
		Domain    string
		Type      string
		Body      []byte
		Sent      bool
		SentAt    time.Time
		CreatedAt time.Time
	}

	TransactionalOutboxOption func(*TransactionalOutbox)
)

func NewTransactionalOutbox(
	domain string,
	eventType string,
	body []byte,
	opts ...TransactionalOutboxOption,
) TransactionalOutbox {
	var to = TransactionalOutbox{
		Domain: domain,
		Type:   eventType,
		Body:   body,
	}

	for _, o := range opts {
		o(&to)
	}

	return to
}

func WithID(id int64) TransactionalOutboxOption {
	return func(to *TransactionalOutbox) {
		to.ID = valueobject.ID(id)
	}
}

func WithCreatedAt(createdAt time.Time) TransactionalOutboxOption {
	return func(to *TransactionalOutbox) {
		to.CreatedAt = createdAt
	}
}
