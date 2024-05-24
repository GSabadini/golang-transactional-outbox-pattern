package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
)

const (
	queryInsertAccount = `INSERT INTO Accounts (Document, CreatedAt) VALUES (?, ?);`
)

type (
	AccountModel struct {
		Document  sql.NullString
		CreatedAt sql.NullTime
	}

	AccountRepository struct {
		db *sql.DB
	}
)

func NewAccountRepository(db *sql.DB) AccountRepository {
	return AccountRepository{db: db}
}

func (a AccountRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	account domain.Account,
) (valueobject.ID, error) {
	ctx, span := opentelemetry.NewSpan(ctx, "repository.account.create")
	defer span.End()

	var model = AccountModel{
		Document:  newNullString(account.Document.String()),
		CreatedAt: newNullTime(account.CreatedAt),
	}

	result, err := tx.ExecContext(
		ctx,
		queryInsertAccount,
		model.Document,
		model.CreatedAt,
	)
	if err != nil {
		opentelemetry.SetError(span, err)
		return valueobject.ID(0), err
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		opentelemetry.SetError(span, err)
		return valueobject.ID(0), err
	}

	return valueobject.ID(accountID), nil
}
