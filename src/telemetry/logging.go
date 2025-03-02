package telemetry

import (
	"context"
	"log/slog"
	"os"

	"github.com/inx51/howlite/resources/config"
	"github.com/phsym/console-slog"
	slogmulti "github.com/samber/slog-multi"
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

	otlpExporter, err := getOtlLogExporter(ctx, conf)
	if err != nil {
		panic(err)
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
		log.WithResource(resource.NewWithAttributes(semconv.SchemaURL)),
	)

	global.SetLoggerProvider(loggerProvider)

	var handlers []slog.Handler

	if conf.DEV_MODE {
		handlers = append(handlers, console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}))
	}

	handlers = append(handlers, otelslog.NewHandler(conf.OTEL_SERVICE_NAME))

	return slog.New(slogmulti.Fanout(handlers...))
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
