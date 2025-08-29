package telemetry

import (
	"context"
	"os"

	"github.com/inx51/howlite-resources/logger"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var loggerProvider *log.LoggerProvider

func newLoggerProvider(ctx context.Context) (*log.LoggerProvider, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	otlpExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		panic(err)
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ProcessPID(os.Getpid()),
			semconv.HostNameKey.String(hostname))),
	)

	return provider, nil
}

func SetupLogging(ctx context.Context) {
	var err error
	loggerProvider, err = newLoggerProvider(ctx)
	if err != nil {
		panic(err)
	}

	global.SetLoggerProvider(loggerProvider)
}

func ShutdownLogging(ctx context.Context) {
	if err := loggerProvider.Shutdown(ctx); err != nil {
		logger.Error(ctx, "Failed to shutdown logger provider for OpenTelemetry", "error", err)
		return
	}
}
