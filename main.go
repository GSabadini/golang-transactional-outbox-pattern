package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/GSabadini/golang-transactional-outbox-pattern/infra"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/aws"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/broker"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/database"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"

	"github.com/shopspring/decimal"
)

func main() {
	var ctx = context.TODO()

	decimal.MarshalJSONWithoutQuotes = false

	tracerShutdown, err := opentelemetry.NewTracer(ctx)
	if err != nil {
		logger.Slog.Error("Tracer connection error", slog.String("error", err.Error()))
		os.Exit(0)
	}
	defer tracerShutdown()

	mysql, mysqlShutdown, err := database.NewMySQL(ctx)
	if err != nil {
		logger.Slog.Error("Database connection error", slog.String("error", err.Error()))
		os.Exit(0)
	}
	defer mysqlShutdown()

	awsConfig, err := aws.NewConfig(ctx)
	if err != nil {
		logger.Slog.Error("AWS connection error", slog.String("error", err.Error()))
		os.Exit(0)
	}

	var dependencies = infra.Dependencies{
		SNS:   broker.NewSNS(awsConfig),
		MySQL: mysql,
	}

	infra.NewHTTPServer().Start(ctx, dependencies)
	infra.NewCron().Start(ctx, dependencies)
}
