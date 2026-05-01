package event

import (
	"context"

	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/tracer"
	"github.com/zeromq/goczmq"
)

type Publisher struct {
	socket *goczmq.Sock
}

func NewPublisher(ctx context.Context, endpoint string) Publisher {
	ctx, span := tracer.StartInfoSpan(ctx, "zeromq.publisher.init")
	defer tracer.SafeEndSpan(span)

	logger.Debug(ctx, "Establishing connection to zero mq publisher", "endpoint", endpoint)
	sock, err := goczmq.NewPub(endpoint)
	if err != nil {
		tracer.SafeRecordError(span, err)
		logger.Error(ctx, "Failed to establish connection to zero mq publisher", "endpoint", endpoint, "error", err)
		return Publisher{}
	}

	logger.Info(ctx, "Zero mq publisher initialized", "endpoint", endpoint)
	return Publisher{
		socket: sock,
	}
}

func (publisher *Publisher) Publish(ctx context.Context, event []byte) {
	ctx, span := tracer.StartDebugSpan(ctx, "zeromq.sendframe")
	defer tracer.SafeEndSpan(span)

	logger.Debug(ctx, "Sending event frame via zero mq", "payload", string(event))
	err := publisher.socket.SendFrame(event, goczmq.FlagNone)
	if err != nil {
		tracer.SafeRecordError(span, err)
		logger.Error(ctx, "Failed to send event frame via zero mq", "error", err)
		return
	}
	logger.Info(ctx, "Event published")
}

func (publisher *Publisher) Stop() {
	publisher.socket.Destroy()
}
