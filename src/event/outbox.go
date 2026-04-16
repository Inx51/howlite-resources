package event

import (
	"context"
	"sync"

	"github.com/inx51/howlite-resources/logger"
	"github.com/vmihailenco/msgpack/v5"
)

type Outbox struct {
	mutex sync.Mutex
	items [][]byte
}

func NewOutbox() Outbox {
	return Outbox{}
}

func (outbox *Outbox) Enqueue(ctx context.Context, event any) {
	msg, err := msgpack.Marshal(event)
	if err != nil {
		logger.Error(ctx, "failed to marshall event", "error", err)
	}

	outbox.mutex.Lock()
	defer outbox.mutex.Unlock()
	outbox.items = append(outbox.items, msg)
}

func (outbox *Outbox) Dequeue(ctx context.Context) []byte {
	outbox.mutex.Lock()
	defer outbox.mutex.Unlock()
	if len(outbox.items) == 0 {
		return nil
	}
	item := outbox.items[0]
	outbox.items = outbox.items[1:]
	return item
}
