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

	"github.com/shopspring/decimal"
)

func main() {
	var ctx = context.TODO()

	decimal.MarshalJSONWithoutQuotes = false

	mysql, err := database.NewMySQL()
	if err != nil {
		logger.Slog.Error("Database connection error", slog.String("error", err.Error()))
		os.Exit(0)
	}

	var dependencies = infra.Dependencies{
		Broker: broker.NewSNS(aws.NewConfig()),
		MySQL:  mysql,
	}

	infra.NewHTTPServer().Start(ctx, dependencies)
	infra.NewCron().Start(ctx, dependencies)
}
