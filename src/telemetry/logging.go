package telemetry

import (
	"context"
	"log/slog"
	"os"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryLogger(conf config.OtelConfiguration) *slog.Logger {
	consoleExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	otlpExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		panic(err)
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(consoleExporter)),
		log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OTEL_SERVICE_NAME),
			semconv.ProcessPID(os.Getpid()),
			semconv.HostNameKey.String(hostname))),
	)

	global.SetLoggerProvider(loggerProvider)

	return slog.New(otelslog.NewHandler(conf.OTEL_SERVICE_NAME))
}
