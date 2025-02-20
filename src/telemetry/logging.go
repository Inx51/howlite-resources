package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryLogger(serviceName string) *slog.Logger {
	stdoutExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	otlpExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint("127.0.0.1:4317"), otlploggrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// Set up logger provider.
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(stdoutExporter)),
		log.WithProcessor(log.NewBatchProcessor(otlpExporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	// shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return otelslog.NewLogger(serviceName)
}
