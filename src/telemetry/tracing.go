package telemetry

import (
	"context"
	"os"

	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var tracerProvider *trace.TracerProvider

func newTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {
	otlpExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		logger.Warn(ctx, "Failed to create tracing exporter for OpenTelemetry, skipping OpenTelemetry tracing", "error", err)
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(otlpExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ProcessPID(os.Getpid()),
			semconv.HostNameKey.String(hostname),
		)),
	)
	return provider, nil
}

func SetupTracing(
	ctx context.Context,
	configuration *configuration.Tracing) {
	var err error
	tracerProvider, err = newTracerProvider(ctx)
	if err != nil {
		return
	}
	otel.SetTracerProvider(tracerProvider)
	tracer.SetupTracer(configuration, tracerProvider != nil)
}

func ShutdownTracing(ctx context.Context) {
	if err := tracerProvider.Shutdown(ctx); err != nil {
		logger.Error(ctx, "OpenTelemetry failed to shutdown tracer provider", "error", err)
		return
	}
}
