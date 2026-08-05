package event

import (
	"context"

	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/tracer"
	"github.com/zeromq/goczmq"
)

type Publisher struct {
	socket *goczmq.Sock
	auth   *goczmq.Auth
}

func (publisher Publisher) IsAvailable() bool {
	return publisher.socket != nil
}

func NewPublisher(ctx context.Context, config configuration.ZeroMqConfiguration) Publisher {
	ctx, span := tracer.StartInfoSpan(ctx, "zeromq.publisher.init")
	defer tracer.SafeEndSpan(span)

	logger.Debug(ctx, "Establishing connection to zero mq publisher", "endpoint", config.ENDPOINT)

	sock := goczmq.NewSock(goczmq.Pub)

	var auth *goczmq.Auth
	if config.CURVE.SERVER_CERT_PATH != "" {
		var err error
		auth, err = setupCurve(sock, config.CURVE)
		if err != nil {
			tracer.SafeRecordError(span, err)
			logger.Error(ctx, "Failed to configure CURVE for zero mq publisher", "error", err)
			sock.Destroy()
			return Publisher{}
		}
	}

	if err := sock.Attach(config.ENDPOINT, true); err != nil {
		tracer.SafeRecordError(span, err)
		logger.Error(ctx, "Failed to establish connection to zero mq publisher", "endpoint", config.ENDPOINT, "error", err)
		sock.Destroy()
		if auth != nil {
			auth.Destroy()
		}
		return Publisher{}
	}

	logger.Info(ctx, "Zero mq publisher initialized", "endpoint", config.ENDPOINT)
	return Publisher{
		socket: sock,
		auth:   auth,
	}
}

// setupCurve loads the publisher's CURVE cert onto sock, marks it as a CURVE
// server, and starts an auth actor enforcing the configured client allowlist
// (or CURVE_ALLOW_ANY if none was configured). Must be called before the
// socket is bound.
func setupCurve(sock *goczmq.Sock, curve configuration.ZeroMqCurveConfiguration) (*goczmq.Auth, error) {
	cert, err := goczmq.NewCertFromFile(curve.SERVER_CERT_PATH)
	if err != nil {
		return nil, err
	}

	sock.SetZapDomain("global")
	cert.Apply(sock)
	sock.SetCurveServer(1)

	allowed := goczmq.CurveAllowAny
	if curve.ALLOWED_CLIENTS_PATH != "" {
		allowed = curve.ALLOWED_CLIENTS_PATH
	}

	auth := goczmq.NewAuth()
	if err := auth.Curve(allowed); err != nil {
		auth.Destroy()
		return nil, err
	}

	return auth, nil
}

func (publisher *Publisher) Publish(ctx context.Context, event []byte) {
	if publisher == nil || publisher.socket == nil {
		logger.Error(ctx, "Zero mq publisher is not available")
		return
	}

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
	if publisher == nil || publisher.socket == nil {
		return
	}

	publisher.socket.Destroy()
	if publisher.auth != nil {
		publisher.auth.Destroy()
	}
}
