package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
)

const (
	queryInsertTransactionalOutbox       = `INSERT INTO TransactionalOutbox (Domain, Type, Body, Sent, SentAt, CreatedAt) VALUES (?, ?, ?, ?, ?, ?);`
	queryFindByUnsentTransactionalOutbox = `SELECT ID, Domain, Type, Body FROM TransactionalOutbox WHERE Sent=0 LIMIT 1 FOR UPDATE SKIP LOCKED;`
	queryMarkToSentTransactionalOutbox   = `UPDATE TransactionalOutbox SET Sent=(?), SentAt=(?) WHERE ID=(?);`
)

type TransactionalOutboxModel struct {
	ID        sql.NullInt64
	Domain    sql.NullString
	Type      sql.NullString
	Body      []byte
	Sent      sql.NullBool
	SentAt    sql.NullTime
	CreatedAt sql.NullTime
}

type TransactionalOutboxRepository struct {
	db *sql.DB
}

func NewTransactionalOutboxRepository(db *sql.DB) TransactionalOutboxRepository {
	return TransactionalOutboxRepository{db: db}
}

func (tor TransactionalOutboxRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	transactionalOutbox domain.TransactionalOutbox,
) error {
	var model = TransactionalOutboxModel{
		Domain:    NewNullString(transactionalOutbox.Domain),
		Type:      NewNullString(transactionalOutbox.Type),
		Body:      transactionalOutbox.Body,
		Sent:      NewNullBool(transactionalOutbox.Sent),
		SentAt:    NewNullTime(transactionalOutbox.SentAt),
		CreatedAt: NewNullTime(transactionalOutbox.CreatedAt),
	}

	_, err := tx.ExecContext(
		ctx,
		queryInsertTransactionalOutbox,
		model.Domain,
		model.Type,
		model.Body,
		model.Sent,
		model.SentAt,
		model.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (tor TransactionalOutboxRepository) FindByUnsent(ctx context.Context) (domain.TransactionalOutbox, error) {
	var model TransactionalOutboxModel

	err := tor.db.QueryRowContext(ctx, queryFindByUnsentTransactionalOutbox).Scan(
		&model.ID,
		&model.Domain,
		&model.Type,
		&model.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TransactionalOutbox{}, nil
		}

		return domain.TransactionalOutbox{}, err
	}

	var transactionalOutbox = domain.NewTransactionalOutbox(
		model.Domain.String,
		model.Type.String,
		model.Body,
		domain.WithID(model.ID.Int64),
	)

	return transactionalOutbox, nil
}

func (tor TransactionalOutboxRepository) MarkToSent(ctx context.Context, id valueobject.ID) error {
	_, err := tor.db.ExecContext(
		ctx,
		queryMarkToSentTransactionalOutbox,
		true,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
