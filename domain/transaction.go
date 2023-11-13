package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"time"

	"github.com/shopspring/decimal"
)

type (
	TransactionCreator interface {
		Create(context.Context, *sql.Tx, Transaction) (valueobject.ID, error)
	}
)

type Transaction struct {
	ID            valueobject.ID            `json:"id"`
	AccountID     valueobject.ID            `json:"account_id"`
	Amount        decimal.Decimal           `json:"amount"`
	Currency      string                    `json:"currency"`
	OperationType valueobject.OperationType `json:"operation_type"`
	CreatedAt     time.Time                 `json:"created_at"`
}

func NewTransaction(
	accountID valueobject.ID,
	amount decimal.Decimal,
	currency string,
	operationType valueobject.OperationType,
	createdAt time.Time,
) Transaction {
	return Transaction{
		AccountID:     accountID,
		Amount:        amount,
		Currency:      currency,
		OperationType: operationType,
		CreatedAt:     createdAt,
	}
}

func (t *Transaction) WithID(id valueobject.ID) {
	t.ID = id
}

func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}
