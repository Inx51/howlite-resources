package event

import (
	"context"
	"encoding/json"

	"github.com/inx51/howlite-resources/logger"
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

func (bus *Bus) Publish(ctx context.Context, eventType string, eventData any) {
	envelope, err := NewEnvelope(eventType, eventData)
	if err != nil {
		logger.Error(ctx, "failed to build envelope", "error", err)
		return
	}

	msg, err := json.Marshal(envelope)
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
