package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
)

type (
	AccountUseCase interface {
		Create(context.Context, CreateAccountInput) (CreateAccountOutput, error)
	}

	CreateAccountInput struct {
		Document valueobject.Document `json:"document"`
	}

	CreateAccountOutput struct {
		ID int64 `json:"id"`
	}

	AccountOrchestrator struct {
		atomic                        domain.Atomic
		accountRepository             domain.AccountRepository
		transactionalOutboxRepository domain.TransactionalOutboxRepository
	}
)

func NewAccountOrchestrator(
	atomic domain.Atomic,
	accountRepository domain.AccountRepository,
	transactionalOutboxRepository domain.TransactionalOutboxRepository,
) AccountOrchestrator {
	return AccountOrchestrator{
		atomic:                        atomic,
		accountRepository:             accountRepository,
		transactionalOutboxRepository: transactionalOutboxRepository,
	}
}

func (a AccountOrchestrator) Create(
	ctx context.Context,
	input CreateAccountInput,
) (CreateAccountOutput, error) {
	ctx, span := opentelemetry.NewSpan(ctx, "usecase.account.execute")
	defer span.End()

	var account = domain.NewAccount(input.Document, time.Now().UTC())

	tx, err := a.atomic.BeginTx(ctx)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateAccountOutput{}, err
	}
	defer a.atomic.Rollback(tx)

	id, err := a.accountRepository.Create(ctx, tx, account)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateAccountOutput{}, err
	}
	account.WithID(id)

	eventBody, err := account.EventBody()
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateAccountOutput{}, err
	}

	err = a.transactionalOutboxRepository.Create(
		ctx,
		tx,
		domain.NewTransactionalOutbox(
			domain.AccountEventDomain,
			domain.AccountEventType,
			eventBody,
			domain.WithCreatedAt(time.Now().UTC()),
		),
	)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateAccountOutput{}, err
	}

	err = a.atomic.Commit(tx)
	if err != nil {
		opentelemetry.SetError(span, err)
		return CreateAccountOutput{}, err
	}

	return CreateAccountOutput{ID: account.ID.Int64()}, nil
}
