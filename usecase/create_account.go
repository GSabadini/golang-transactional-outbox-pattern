package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain/valueobject"
)

type CreateAccountUseCase interface {
	Execute(context.Context, CreateAccountInput) (CreateAccountOutput, error)
}

type CreateAccountInput struct {
	Document valueobject.Document `json:"Document"`
}

type CreateAccountOutput struct {
	ID int64 `json:"id"`
}

type CreateAccountOrchestrate struct {
	atomic                        domain.Atomic
	accountRepository             domain.AccountRepository
	transactionalOutboxRepository domain.TransactionalOutboxRepository
}

func NewCreateAccountOrchestrate(
	atomic domain.Atomic,
	accountRepository domain.AccountRepository,
	transactionalOutboxRepository domain.TransactionalOutboxRepository,
) CreateAccountOrchestrate {
	return CreateAccountOrchestrate{
		atomic:                        atomic,
		accountRepository:             accountRepository,
		transactionalOutboxRepository: transactionalOutboxRepository,
	}
}

func (cao CreateAccountOrchestrate) Execute(
	ctx context.Context,
	input CreateAccountInput,
) (CreateAccountOutput, error) {
	accountID, err := cao.performTransactionalOperation(ctx, domain.NewAccount(input.Document, time.Now().UTC()))
	if err != nil {
		return CreateAccountOutput{}, err
	}

	return CreateAccountOutput{ID: accountID.Int64()}, nil
}

func (cao CreateAccountOrchestrate) performTransactionalOperation(
	ctx context.Context,
	account domain.Account,
) (valueobject.ID, error) {
	tx, err := cao.atomic.BeginTx(ctx)
	if err != nil {
		return valueobject.ID(0), err
	}

	defer cao.atomic.Rollback(tx)

	id, err := cao.accountRepository.Create(ctx, tx, account)
	if err != nil {
		return valueobject.ID(0), err
	}
	account.WithID(id)

	eventBody, err := account.ToJSON()
	if err != nil {
		return valueobject.ID(0), err
	}

	err = cao.transactionalOutboxRepository.Create(
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
		return valueobject.ID(0), err
	}

	err = tx.Commit()
	if err != nil {
		return valueobject.ID(0), err
	}

	return id, nil
}
