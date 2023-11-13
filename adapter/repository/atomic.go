package repository

import (
	"context"
	"database/sql"
)

type Atomic struct {
	db *sql.DB
}

func NewAtomic(db *sql.DB) Atomic {
	return Atomic{db: db}
}

func (a Atomic) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return a.db.BeginTx(ctx, &sql.TxOptions{})
}

func (a Atomic) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (a Atomic) Rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}
