package domain

import (
	"context"
	"database/sql"
)

type (
	Atomic interface {
		BeginTx(context.Context) (*sql.Tx, error)
		Commit(*sql.Tx) error
		Rollback(*sql.Tx)
	}
)
