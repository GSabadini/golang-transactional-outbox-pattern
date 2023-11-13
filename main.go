package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/handler"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/background"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/producer"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/repository"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/database"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/server"
	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"

	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/middleware"
	"github.com/shopspring/decimal"
)

func main() {
	mysql, err := database.NewMySQL()
	if err != nil {
		logger.Slog.Error("Database connection error", slog.String("error", err.Error()))
		os.Exit(0)
	}

	decimal.MarshalJSONWithoutQuotes = false

	var (
		ctx = context.TODO()

		echoServer = server.NewEcho()

		cronScheduler = gocron.NewScheduler(time.UTC)

		publish                       = producer.NewPublish()
		atomic                        = repository.NewAtomic(mysql)
		transactionRepository         = repository.NewTransactionRepository(mysql)
		transactionalOutboxRepository = repository.NewTransactionalOutboxRepository(mysql)

		transactionUseCase = usecase.NewCreateTransactionOrchestrate(
			atomic,
			transactionRepository,
			transactionalOutboxRepository,
		)

		transactionHandler = handler.NewCreateTransactionHandler(transactionUseCase)

		transactionalOutboxBackground = background.NewTransactionalOutbox(
			publish,
			transactionalOutboxRepository,
			transactionalOutboxRepository,
		)
	)

	echoServer.Use(middleware.Recover())

	echoServer.POST("/transaction", transactionHandler.Handle)

	go func() {
		echoServer.Logger.Fatal(echoServer.Start(":8080"))
	}()

	scheduleCron(ctx, cronScheduler, transactionalOutboxBackground)

	cronScheduler.StartBlocking()
	cronScheduler.Stop()
}

func scheduleCron(ctx context.Context, scheduler *gocron.Scheduler, transactionalOutbox background.TransactionalOutbox) {
	_, err := scheduler.Every(5).Seconds().Do(func() {
		err := transactionalOutbox.Process(ctx)
		if err != nil {
			logger.Slog.Error("Error execute the outbox process", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		logger.Slog.Error("Error execute the job", slog.String("error", err.Error()))
		os.Exit(0)
	}
}
