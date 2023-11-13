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

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return TransactionRepository{db: db}
}

func (tr TransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction domain.Transaction) (valueobject.ID, error) {
	result, err := tx.ExecContext(
		ctx,
		queryInsertTransaction,
		transaction.AccountID,
		transaction.Amount,
		transaction.Currency,
		transaction.OperationType,
		transaction.CreatedAt,
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
