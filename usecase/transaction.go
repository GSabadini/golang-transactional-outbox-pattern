package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"

	"github.com/shopspring/decimal"
)

type (
	TransactionUseCase interface {
		Create(context.Context, CreateTransactionInput) (CreateTransactionOutput, error)
	}

	CreateTransactionInput struct {
		AccountID     valueobject.ID            `json:"account_id"`
		Currency      valueobject.Currency      `json:"currency"`
		OperationType valueobject.OperationType `json:"operation_type"`
		Amount        decimal.Decimal           `json:"amount"`
	}

	CreateTransactionOutput struct {
		ID int64 `json:"id"`
	}

	TransactionOrchestrator struct {
		atomic                        domain.Atomic
		transactionRepository         domain.TransactionRepository
		transactionalOutboxRepository domain.TransactionalOutboxRepository
	}
)

func NewTransactionOrchestrator(
	atomic domain.Atomic,
	transactionRepository domain.TransactionRepository,
	transactionalOutboxRepository domain.TransactionalOutboxRepository,
) TransactionOrchestrator {
	return TransactionOrchestrator{
		atomic:                        atomic,
		transactionRepository:         transactionRepository,
		transactionalOutboxRepository: transactionalOutboxRepository,
	}
}

func (t TransactionOrchestrator) Create(
	ctx context.Context,
	input CreateTransactionInput,
) (CreateTransactionOutput, error) {
	ctx, span := opentelemetry.NewSpan(ctx, "usecase.transaction.create")
	defer span.End()

	var transaction = domain.NewTransaction(
		input.AccountID,
		input.Amount,
		input.Currency,
		input.OperationType,
		time.Now().UTC(),
	)

	tx, err := t.atomic.BeginTx(ctx)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateTransactionOutput{}, err
	}
	defer t.atomic.Rollback(tx)

	id, err := t.transactionRepository.Create(ctx, tx, transaction)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateTransactionOutput{}, err
	}
	transaction.WithID(id)

	eventBody, err := transaction.EventBody()
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateTransactionOutput{}, err
	}

	err = t.transactionalOutboxRepository.Create(
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
		opentelemetry.SetError(span, err)
		return CreateTransactionOutput{}, err
	}

	err = t.atomic.Commit(tx)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateTransactionOutput{}, err
	}

	return CreateTransactionOutput{ID: transaction.ID.Int64()}, nil
}
