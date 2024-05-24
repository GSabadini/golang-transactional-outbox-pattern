package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"

	"github.com/shopspring/decimal"
)

const (
	TransactionEventDomain = "transactions"
	TransactionEventType   = "transaction-created"
)

type (
	TransactionRepository interface {
		Create(context.Context, *sql.Tx, Transaction) (valueobject.ID, error)
	}

	Transaction struct {
		ID            valueobject.ID            `json:"id"`
		AccountID     valueobject.ID            `json:"account_id"`
		Currency      valueobject.Currency      `json:"currency"`
		OperationType valueobject.OperationType `json:"operation_type"`
		Amount        decimal.Decimal           `json:"amount"`
		CreatedAt     time.Time                 `json:"created_at"`
	}
)

func NewTransaction(
	accountID valueobject.ID,
	amount decimal.Decimal,
	currency valueobject.Currency,
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

func (t *Transaction) EventBody() ([]byte, error) {
	return json.Marshal(t)
}
