package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
)

const (
	queryInsertAccount = `INSERT INTO Accounts (Document, CreatedAt) VALUES (?, ?);`
)

type AccountModel struct {
	Document  sql.NullString
	CreatedAt sql.NullTime
}

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return AccountRepository{db: db}
}

func (ar AccountRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	account domain.Account,
) (valueobject.ID, error) {
	var model = AccountModel{
		Document:  NewNullString(account.Document.String()),
		CreatedAt: NewNullTime(account.CreatedAt),
	}

	result, err := tx.ExecContext(
		ctx,
		queryInsertAccount,
		model.Document,
		model.CreatedAt,
	)
	if err != nil {
		return valueobject.ID(0), err
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		return valueobject.ID(0), err
	}

	return valueobject.ID(accountID), nil
}
