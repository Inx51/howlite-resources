package telemetry

import (
	"context"
	"log/slog"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryLogger(conf config.OtelConfiguration) *slog.Logger {
	stdoutExporter, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	otlpExporter, err := getOtlLogExporter(ctx, conf)
	if err != nil {
		panic(err)
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(stdoutExporter)),
		log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
		log.WithResource(resource.NewWithAttributes(semconv.SchemaURL)),
	)

	global.SetLoggerProvider(loggerProvider)

	return otelslog.NewLogger(conf.OTEL_SERVICE_NAME)
}

func getOtlLogExporter(ctx context.Context, conf config.OtelConfiguration) (log.Exporter, error) {
	protocol := conf.OTEL_EXPORTER_OTLP_PROTOCOL
	if conf.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL != "" {
		protocol = conf.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL
	}

	if protocol == "http" {
		return otlploghttp.New(ctx)
	}

	return otlploggrpc.New(ctx)
}
