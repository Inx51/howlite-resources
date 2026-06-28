package event

import (
	"context"
	"time"

	"github.com/inx51/howlite-resources/logger"
)

type OutboxWorker struct {
	outbox    *Outbox
	publisher *Publisher
	ticker    *time.Ticker
}

func NewOutboxWorker(ctx context.Context, outbox *Outbox, publisher *Publisher) OutboxWorker {
	return OutboxWorker{
		outbox:    outbox,
		publisher: publisher,
		ticker:    time.NewTicker(10 * time.Millisecond),
	}
}

func (worker *OutboxWorker) Start(ctx context.Context) {
	if worker.ticker == nil {
		logger.Info(ctx, "Background worker for outbox did not start, missing endpoint configuration for publisher")
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "Outbox worker stopped")
			return
		case <-worker.ticker.C:
			message := worker.outbox.Dequeue(ctx)
			if message == nil {
				logger.Debug(ctx, "No new messages in outbox")
				continue
			}

			worker.publisher.Publish(ctx, message)
			logger.Info(ctx, "Published event from outbox")
		}
	}
}

func (worker *OutboxWorker) Stop(ctx context.Context) {
	worker.outbox.Close(ctx)
	worker.ticker.Stop()

}
