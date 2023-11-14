package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"

	"github.com/shopspring/decimal"
)

type CreateTransactionUseCase interface {
	Execute(context.Context, CreateTransactionInput) (CreateTransactionOutput, error)
}

type CreateTransactionInput struct {
	AccountID     valueobject.ID            `json:"account_id"`
	Currency      valueobject.Currency      `json:"currency"`
	OperationType valueobject.OperationType `json:"operation_type"`
	Amount        decimal.Decimal           `json:"amount"`
}

type CreateTransactionOutput struct {
	ID int64 `json:"id"`
}

type CreateTransactionOrchestrate struct {
	atomic                        domain.Atomic
	transactionRepository         domain.TransactionRepository
	transactionalOutboxRepository domain.TransactionalOutboxRepository
}

func NewCreateTransactionOrchestrate(
	atomic domain.Atomic,
	transactionRepository domain.TransactionRepository,
	transactionalOutboxRepository domain.TransactionalOutboxRepository,
) CreateTransactionOrchestrate {
	return CreateTransactionOrchestrate{
		atomic:                        atomic,
		transactionRepository:         transactionRepository,
		transactionalOutboxRepository: transactionalOutboxRepository,
	}
}

func (cto CreateTransactionOrchestrate) Execute(
	ctx context.Context,
	input CreateTransactionInput,
) (CreateTransactionOutput, error) {
	transactionID, err := cto.performTransactionalOperation(ctx, domain.NewTransaction(
		input.AccountID,
		input.Amount,
		input.Currency,
		input.OperationType,
		time.Now().UTC(),
	))
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

	id, err := cto.transactionRepository.Create(ctx, tx, transaction)
	if err != nil {
		return valueobject.ID(0), err
	}
	transaction.WithID(id)

	eventBody, err := transaction.ToJSON()
	if err != nil {
		return valueobject.ID(0), err
	}

	err = cto.transactionalOutboxRepository.Create(
		ctx,
		tx,
		domain.NewTransactionalOutbox(
			domain.TransactionEventDomain,
			domain.TransactionEventType,
			eventBody,
			domain.WithCreatedAt(time.Now().UTC()),
		),
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
