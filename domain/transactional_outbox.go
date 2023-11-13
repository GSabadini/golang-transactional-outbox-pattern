package domain

import (
	"context"
	"database/sql"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"time"
)

type (
	TransactionalOutboxCreator interface {
		Create(context.Context, *sql.Tx, TransactionalOutbox) error
	}

	TransactionalOutboxFinder interface {
		FindByUnsent(context.Context) (TransactionalOutbox, error)
	}

	TransactionalOutboxUpdater interface {
		MarkToSent(context.Context, valueobject.ID) error
	}
)

type TransactionalOutbox struct {
	ID        valueobject.ID
	Target    string
	EventType string
	Body      []byte
	Sent      bool
	CreatedAt time.Time
}

func NewTransactionalOutbox(body []byte, sent bool, createdAt time.Time) TransactionalOutbox {
	return TransactionalOutbox{
		Body:      body,
		Sent:      sent,
		CreatedAt: createdAt,
	}
}
