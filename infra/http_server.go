package infra

import (
	"context"
	"errors"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/handler"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/repository"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/env"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/server"
	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/middleware"
)

type (
	HTTPServer struct{}
)

func NewHTTPServer() HTTPServer {
	return HTTPServer{}
}

func (h HTTPServer) Start(ctx context.Context, dependencies Dependencies) {
	var (
		echoServer = server.NewEcho()
		v1         = echoServer.Group("/v1")
	)

	echoServer.Use(middleware.Recover(), middleware.RequestID())

	var (
		transactionHandler = handler.NewTransactionHandler(
			usecase.NewTransactionOrchestrator(
				repository.NewAtomic(dependencies.MySQL),
				repository.NewTransactionRepository(dependencies.MySQL),
				repository.NewTransactionalOutboxRepository(dependencies.MySQL),
			),
		)

		accountHandler = handler.NewAccountHandler(
			usecase.NewAccountOrchestrator(
				repository.NewAtomic(dependencies.MySQL),
				repository.NewAccountRepository(dependencies.MySQL),
				repository.NewTransactionalOutboxRepository(dependencies.MySQL),
			),
		)
	)

	v1.POST("/transactions", transactionHandler.Create)
	v1.POST("/accounts", accountHandler.Create)

	go func() {
		err := echoServer.Start(env.ServerPort)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			echoServer.Logger.Fatal(err)
		}
	}()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-quit

		ctxTimeout, cancel := context.WithTimeout(ctx, env.ServerShutdownTimeout)
		defer cancel()

		if err := echoServer.Shutdown(ctxTimeout); err != nil {
			echoServer.Logger.Fatal(err)
		}
	}()
}
