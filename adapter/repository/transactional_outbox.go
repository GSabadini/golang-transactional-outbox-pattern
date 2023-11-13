package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
)

const (
	queryInsertTransactionalOutbox       = `INSERT INTO TransactionalOutbox (Body, Sent, CreatedAt) VALUES (?, ?, ?);`
	queryFindByUnsentTransactionalOutbox = `SELECT * FROM TransactionalOutbox WHERE Sent=0 LIMIT 1;`
	queryMarkToSentTransactionalOutbox   = `UPDATE TransactionalOutbox SET Sent=(?) WHERE ID=(?);`
)

type TransactionalOutboxRepository struct {
	db *sql.DB
}

func NewTransactionalOutboxRepository(db *sql.DB) TransactionalOutboxRepository {
	return TransactionalOutboxRepository{db: db}
}

func (tor TransactionalOutboxRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	transactionOutbox domain.TransactionalOutbox,
) error {
	_, err := tx.ExecContext(
		ctx,
		queryInsertTransactionalOutbox,
		transactionOutbox.Body,
		transactionOutbox.Sent,
		transactionOutbox.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (tor TransactionalOutboxRepository) FindByUnsent(ctx context.Context) (domain.TransactionalOutbox, error) {
	var transactionOutbox domain.TransactionalOutbox

	err := tor.db.QueryRowContext(ctx, queryFindByUnsentTransactionalOutbox).Scan(
		&transactionOutbox.ID,
		&transactionOutbox.Body,
		&transactionOutbox.Sent,
		&transactionOutbox.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TransactionalOutbox{}, nil
		}

		return domain.TransactionalOutbox{}, err
	}

	return transactionOutbox, nil
}

func (tor TransactionalOutboxRepository) MarkToSent(ctx context.Context, id valueobject.ID) error {
	_, err := tor.db.ExecContext(ctx, queryMarkToSentTransactionalOutbox, true, id)
	if err != nil {
		return err
	}

	return nil
}
