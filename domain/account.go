package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
)

const (
	AccountEventDomain = "accounts"
	AccountEventType   = "account-created"
)

type (
	AccountRepository interface {
		Create(context.Context, *sql.Tx, Account) (valueobject.ID, error)
	}

	Account struct {
		ID        valueobject.ID       `json:"id"`
		Document  valueobject.Document `json:"document"`
		CreatedAt time.Time            `json:"created_at"`
	}
)

func NewAccount(document valueobject.Document, createdAt time.Time) Account {
	return Account{
		Document:  document,
		CreatedAt: createdAt,
	}
}

func (a *Account) WithID(id valueobject.ID) {
	a.ID = id
}

func (a *Account) EventBody() ([]byte, error) {
	return json.Marshal(a)
}
