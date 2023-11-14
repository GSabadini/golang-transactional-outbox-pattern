package infra

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/handler"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/repository"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/server"
	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"

	"github.com/labstack/echo/middleware"
)

type HTTPServer struct{}

func NewHTTPServer() HTTPServer {
	return HTTPServer{}
}

func (h HTTPServer) Start(ctx context.Context, dependencies Dependencies) {
	var (
		echoServer = server.NewEcho()
		v1         = echoServer.Group("/v1")
	)

	echoServer.Use(middleware.Recover())

	v1.POST("/transactions", buildCreateTransactionHandler(dependencies).Handle)
	v1.POST("/accounts", buildCreteAccountHandler(dependencies).Handle)

	go func() {
		echoServer.Logger.Fatal(echoServer.Start(":8080"))
	}()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		if err := echoServer.Shutdown(ctxWithTimeout); err != nil {
			echoServer.Logger.Fatal(err)
		}
	}()
}

func buildCreateTransactionHandler(dependencies Dependencies) handler.CreateTransactionHandler {
	return handler.NewCreateTransactionHandler(
		usecase.NewCreateTransactionOrchestrate(
			repository.NewAtomic(dependencies.MySQL),
			repository.NewTransactionRepository(dependencies.MySQL),
			repository.NewTransactionalOutboxRepository(dependencies.MySQL),
		),
	)
}

func buildCreteAccountHandler(dependencies Dependencies) handler.CreateAccountHandler {
	return handler.NewCreateAccountHandler(
		usecase.NewCreateAccountOrchestrate(
			repository.NewAtomic(dependencies.MySQL),
			repository.NewAccountRepository(dependencies.MySQL),
			repository.NewTransactionalOutboxRepository(dependencies.MySQL),
		),
	)
}
