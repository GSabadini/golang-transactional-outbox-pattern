package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
)

const (
	queryInsertTransaction = `INSERT INTO Transactions (Account_ID, Amount, Currency, OperationType, CreatedAt) VALUES (?, ?, ?, ?, ?);`
)

type TransactionModel struct {
	AccountID     sql.NullInt64
	Currency      sql.NullString
	OperationType sql.NullString
	Amount        sql.NullFloat64
	CreatedAt     sql.NullTime
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return TransactionRepository{db: db}
}

func (tr TransactionRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	transaction domain.Transaction,
) (valueobject.ID, error) {
	var model = TransactionModel{
		AccountID:     NewNullInt64(transaction.AccountID.Int64()),
		Currency:      NewNullString(transaction.Currency.String()),
		OperationType: NewNullString(transaction.OperationType.String()),
		Amount:        NewNullFloat64(transaction.Amount.InexactFloat64()),
		CreatedAt:     NewNullTime(transaction.CreatedAt),
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
		return valueobject.ID(0), err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return valueobject.ID(0), err
	}

	return valueobject.ID(transactionID), nil
}
