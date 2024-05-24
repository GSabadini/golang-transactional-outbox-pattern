package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
)

const (
	queryInsertTransaction = `INSERT INTO Transactions (Account_ID, Amount, Currency, OperationType, CreatedAt) VALUES (?, ?, ?, ?, ?);`
)

type (
	TransactionModel struct {
		AccountID     sql.NullInt64
		Currency      sql.NullString
		OperationType sql.NullString
		Amount        sql.NullFloat64
		CreatedAt     sql.NullTime
	}

	TransactionRepository struct {
		db *sql.DB
	}
)

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return TransactionRepository{db: db}
}

func (t TransactionRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	transaction domain.Transaction,
) (valueobject.ID, error) {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.transaction.create")
	defer span.End()

	var model = TransactionModel{
		AccountID:     newNullInt64(transaction.AccountID.Int64()),
		Currency:      newNullString(transaction.Currency.String()),
		OperationType: newNullString(transaction.OperationType.String()),
		Amount:        newNullFloat64(transaction.Amount.InexactFloat64()),
		CreatedAt:     newNullTime(transaction.CreatedAt),
	}

	result, err := tx.ExecContext(
		ctx,
		queryInsertTransaction,
		model.AccountID,
		model.Amount,
		model.Currency,
		model.OperationType,
		model.CreatedAt,
	)
	if err != nil {
		opentelemetry.SetError(span, err)
		return valueobject.ID(0), err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		opentelemetry.SetError(span, err)
		return valueobject.ID(0), err
	}

	return valueobject.ID(transactionID), nil
}
