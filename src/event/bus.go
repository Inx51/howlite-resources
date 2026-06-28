package event

import (
	"context"

	"github.com/inx51/howlite-resources/logger"
	"github.com/vmihailenco/msgpack/v5"
)

type Bus struct {
	outbox    *Outbox
	publisher *Publisher
}

func NewBus(publisher *Publisher, outbox *Outbox) *Bus {
	return &Bus{
		publisher: publisher,
		outbox:    outbox,
	}
}

func (bus *Bus) Publish(ctx context.Context, event any) {
	msg, err := msgpack.Marshal(event)
	if err != nil {
		logger.Error(ctx, "failed to marshall event", "error", err)
		return
	}

	logger.Debug(ctx, "Sending event", "event", string(msg))
	if bus.outbox != nil {
		bus.outbox.Enqueue(ctx, msg)
		return
	}

	if bus.publisher != nil {
		bus.publisher.Publish(ctx, msg)
	}
}
