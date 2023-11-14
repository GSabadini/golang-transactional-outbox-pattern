package infra

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/background"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/producer"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/repository"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"

	"github.com/go-co-op/gocron"
)

type Cron struct{}

func NewCron() Cron {
	return Cron{}
}

func (c Cron) Start(ctx context.Context, dependencies Dependencies) {
	var cronScheduler = gocron.NewScheduler(time.UTC)

	c.runTransactionalOutboxBackground(
		ctx,
		cronScheduler,
		buildTransactionalOutboxBackground(dependencies),
	)

	cronScheduler.StartBlocking()
	cronScheduler.Stop()
	cronScheduler.Clear()
}

func (c Cron) runTransactionalOutboxBackground(
	ctx context.Context,
	cronScheduler *gocron.Scheduler,
	transactionalOutboxBackground background.TransactionalOutbox,
) {
	_, err := cronScheduler.Every(5).Seconds().Do(func() {
		err := transactionalOutboxBackground.Process(ctx)
		if err != nil {
			logger.Slog.Error("Error execute the outbox process", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		logger.Slog.Error("Error execute the job", slog.String("error", err.Error()))
		os.Exit(0)
	}
}

func buildTransactionalOutboxBackground(dependencies Dependencies) background.TransactionalOutbox {
	return background.NewTransactionalOutbox(
		producer.NewProducer(dependencies.Broker),
		repository.NewTransactionalOutboxRepository(dependencies.MySQL),
	)
}
