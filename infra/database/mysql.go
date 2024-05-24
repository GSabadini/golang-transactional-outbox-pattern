package database

import (
	"context"
	"database/sql"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/env"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
)

func NewMySQL(ctx context.Context) (*sql.DB, func(), error) {
	driveName, err := otelsql.Register(
		env.DBDriver,
		otelsql.WithAttributes(
			semconv.DBSystemMySQL,
			semconv.DBNameKey.String(env.DBName),
		),
		otelsql.WithSQLCommenter(true),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			Ping:           true,
			DisableErrSkip: true,
			RecordError:    func(err error) bool { return true },
		}),
	)
	if err != nil {
		return &sql.DB{}, nil, err
	}

	config := &mysql.Config{
		User:              env.DBUser,
		Passwd:            env.DBPassword,
		Addr:              env.DBEndpoint,
		DBName:            env.DBName,
		Net:               "tcp",
		Timeout:           10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
		CheckConnLiveness: true,
		ParseTime:         true,
	}

	db, err := sql.Open(driveName, config.FormatDSN())
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	if err = db.PingContext(ctx); err != nil {
		return nil, nil, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
		semconv.DBNameKey.String(env.DBName),
	))
	if err != nil {
		return nil, nil, err
	}

	shutdown := func() {
		err := db.Close()
		if err != nil {
			return
		}
	}

	return db, shutdown, nil
}
