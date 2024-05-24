package infra

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/env"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/background"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/producer"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/repository"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"

	"github.com/go-co-op/gocron"
)

type (
	Cron struct{}
)

func NewCron() Cron {
	return Cron{}
}

func (c Cron) Start(ctx context.Context, dependencies Dependencies) {
	var (
		cronScheduler = gocron.NewScheduler(time.UTC)

		transactionalOutbox = background.NewTransactionalOutbox(
			producer.NewProducer(dependencies.SNS, env.SNSEventTopic),
			repository.NewTransactionalOutboxRepository(dependencies.MySQL),
		)
	)

	c.runTransactionalOutboxBackground(
		ctx,
		cronScheduler,
		env.CronInterval,
		transactionalOutbox,
	)

	cronScheduler.StartBlocking()
	cronScheduler.Stop()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-quit

		cronScheduler.Clear()
	}()
}

func (c Cron) runTransactionalOutboxBackground(
	ctx context.Context,
	cronScheduler *gocron.Scheduler,
	interval time.Duration,
	transactionalOutboxBackground background.TransactionalOutbox,
) {
	_, err := cronScheduler.Every(interval).Do(func() {
		err := transactionalOutboxBackground.Process(ctx)
		if err != nil {
			logger.Slog.Error("Error execute the transactional outbox process", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		logger.Slog.Error("Error execute the transactional outbox job", slog.String("error", err.Error()))
	}
}
