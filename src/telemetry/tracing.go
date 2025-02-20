package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryTracing(serviceName string) *trace.TracerProvider {
	// shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	consoleExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	otlpExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("127.0.0.1:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(consoleExporter),
		trace.WithBatcher(otlpExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp
}
