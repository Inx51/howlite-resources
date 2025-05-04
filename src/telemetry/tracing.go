package telemetry

import (
	"context"
	"os"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryTracer(conf config.OtelConfiguration) *trace.TracerProvider {
	ctx := context.Background()
	otlpExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(otlpExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OTEL_SERVICE_NAME),
			semconv.ProcessPID(os.Getpid()),
			semconv.HostNameKey.String(hostname),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider
}
