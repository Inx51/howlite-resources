package telemetry

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryLogger(conf config.OtelConfiguration) *slog.Logger {
	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	otlpExporter, err := getOtlpLogExporter(ctx, conf)
	if err != nil {
		panic(err)
	}

	loggerProvider := log.NewLoggerProvider(
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

func getOtlpLogExporter(ctx context.Context, conf config.OtelConfiguration) (log.Exporter, error) {
	protocol := conf.OTEL_EXPORTER_OTLP_PROTOCOL
	if conf.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL != "" {
		protocol = conf.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL
	}

	if strings.HasPrefix(protocol, "http") {
		return otlploghttp.New(ctx)
	}

	return otlploggrpc.New(ctx)
}
