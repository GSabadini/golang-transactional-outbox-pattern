package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
)

const (
	queryInsertTransactionalOutbox          = `INSERT INTO TransactionalOutbox (Domain, Type, Body, Sent, SentAt, CreatedAt) VALUES (?, ?, ?, ?, ?, ?);`
	queryListByUnsentTransactionalOutbox    = `SELECT ID, Domain, Type, Body FROM TransactionalOutbox WHERE Sent=0 LIMIT 10 FOR UPDATE SKIP LOCKED;`
	queryMarkToSentTransactionalOutbox      = `UPDATE TransactionalOutbox SET Sent=(?), SentAt=(?) WHERE ID=(?);`
	queryMarkToSentTransactionalOutboxBatch = `UPDATE TransactionalOutbox SET Sent=(?), SentAt=(?) WHERE ID IN (%s);`
)

type (
	TransactionalOutboxModel struct {
		ID        sql.NullInt64
		Domain    sql.NullString
		Type      sql.NullString
		Body      []byte
		Sent      sql.NullBool
		SentAt    sql.NullTime
		CreatedAt sql.NullTime
	}

	TransactionalOutboxRepository struct {
		db *sql.DB
	}
)

func NewTransactionalOutboxRepository(db *sql.DB) TransactionalOutboxRepository {
	return TransactionalOutboxRepository{db: db}
}

func (t TransactionalOutboxRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	transactionalOutbox domain.TransactionalOutbox,
) error {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.transactional_outbox.create")
	defer span.End()

	var model = TransactionalOutboxModel{
		Domain:    newNullString(transactionalOutbox.Domain),
		Type:      newNullString(transactionalOutbox.Type),
		Body:      transactionalOutbox.Body,
		Sent:      newNullBool(transactionalOutbox.Sent),
		SentAt:    newNullTime(transactionalOutbox.SentAt),
		CreatedAt: newNullTime(transactionalOutbox.CreatedAt),
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
		opentelemetry.SetError(span, err)
		return err
	}

	return nil
}

func (t TransactionalOutboxRepository) ListByUnsent(ctx context.Context) ([]domain.TransactionalOutbox, error) {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.transactional_outbox.list_by_unsent")
	defer span.End()

	var (
		model                 TransactionalOutboxModel
		transactionOutboxList []domain.TransactionalOutbox
	)

	rows, err := t.db.QueryContext(ctx, queryListByUnsentTransactionalOutbox)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.TransactionalOutbox{}, nil
		}

		opentelemetry.SetError(span, err)
		return []domain.TransactionalOutbox{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&model.ID, &model.Domain, &model.Type, &model.Body)
		if err != nil {
			return []domain.TransactionalOutbox{}, err
		}

		transactionOutboxList = append(transactionOutboxList, domain.NewTransactionalOutbox(
			model.Domain.String,
			model.Type.String,
			model.Body,
			domain.WithID(model.ID.Int64),
		))
	}

	err = rows.Close()
	if err != nil {
		return []domain.TransactionalOutbox{}, err
	}

	if err = rows.Err(); err != nil {
		return []domain.TransactionalOutbox{}, err
	}

	return transactionOutboxList, nil
}

func (t TransactionalOutboxRepository) MarkToSent(ctx context.Context, id valueobject.ID) error {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.transactional_outbox.mark_to_sent")
	defer span.End()

	_, err := t.db.ExecContext(
		ctx,
		queryMarkToSentTransactionalOutbox,
		true,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		opentelemetry.SetError(span, err)
		return err
	}

	return nil
}

func (t TransactionalOutboxRepository) MarkToSentBatch(ctx context.Context, idList []valueobject.ID) error {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.transactional_outbox.mark_to_sent_batch")
	defer span.End()

	var idListStr []string

	for _, id := range idList {
		idListStr = append(idListStr, id.String())
	}

	_, err := t.db.ExecContext(
		ctx,
		fmt.Sprintf(queryMarkToSentTransactionalOutboxBatch, strings.Join(idListStr, ",")),
		true,
		time.Now().UTC(),
	)
	if err != nil {
		opentelemetry.SetError(span, err)
		return err
	}

	return nil
}
