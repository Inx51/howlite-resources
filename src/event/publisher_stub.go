//go:build windows || !cgo

package event

import (
	"context"

	"github.com/inx51/howlite-resources/logger"
)

type Publisher struct{}

func (publisher Publisher) IsAvailable() bool {
	return false
}

func NewPublisher(ctx context.Context, endpoint string) Publisher {
	logger.Error(ctx, "Zero mq publisher is unavailable in this build", "endpoint", endpoint)
	return Publisher{}
}

func (publisher *Publisher) Publish(ctx context.Context, event []byte) {
	logger.Error(ctx, "Skipping event publish because Zero mq support is unavailable", "payload_size", len(event))
}

func (publisher *Publisher) Stop() {}