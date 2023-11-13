package usecase

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"

	"github.com/shopspring/decimal"
)

type CreateTransactionUseCase interface {
	Execute(context.Context, CreateTransactionInput) (CreateTransactionOutput, error)
}

type CreateTransactionInput struct {
	AccountID     valueobject.ID            `json:"account_id"`
	Amount        decimal.Decimal           `json:"amount"`
	Currency      string                    `json:"currency"`
	OperationType valueobject.OperationType `json:"operation_type"`
}

type CreateTransactionOutput struct {
	ID int64 `json:"id"`
}

type CreateTransactionOrchestrate struct {
	atomic                     domain.Atomic
	transactionCreator         domain.TransactionCreator
	transactionalOutboxCreator domain.TransactionalOutboxCreator
}

func NewCreateTransactionOrchestrate(
	atomic domain.Atomic,
	transactionCreator domain.TransactionCreator,
	transactionalOutboxCreator domain.TransactionalOutboxCreator,
) CreateTransactionOrchestrate {
	return CreateTransactionOrchestrate{
		atomic:                     atomic,
		transactionCreator:         transactionCreator,
		transactionalOutboxCreator: transactionalOutboxCreator,
	}
}

func (cto CreateTransactionOrchestrate) Execute(
	ctx context.Context,
	input CreateTransactionInput,
) (CreateTransactionOutput, error) {
	transaction := domain.NewTransaction(
		input.AccountID,
		input.Amount,
		input.Currency,
		input.OperationType,
		time.Now().UTC(),
	)

	transactionID, err := cto.performTransactionalOperation(ctx, transaction)
	if err != nil {
		return CreateTransactionOutput{}, err
	}

	return CreateTransactionOutput{ID: transactionID.Int64()}, nil
}

func (cto CreateTransactionOrchestrate) performTransactionalOperation(
	ctx context.Context,
	transaction domain.Transaction,
) (valueobject.ID, error) {
	tx, err := cto.atomic.BeginTx(ctx)
	if err != nil {
		return valueobject.ID(0), err
	}

	defer cto.atomic.Rollback(tx)

	id, err := cto.transactionCreator.Create(ctx, tx, transaction)
	if err != nil {
		return valueobject.ID(0), err
	}
	transaction.WithID(id)

	transactionJSON, err := transaction.ToJSON()
	if err != nil {
		return valueobject.ID(0), err
	}

	err = cto.transactionalOutboxCreator.Create(
		ctx,
		tx,
		domain.NewTransactionalOutbox(transactionJSON, false, time.Now().UTC()),
	)
	if err != nil {
		return valueobject.ID(0), err
	}

	err = tx.Commit()
	if err != nil {
		return valueobject.ID(0), err
	}

	return id, nil
}
