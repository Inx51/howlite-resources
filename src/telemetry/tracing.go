package telemetry

import (
	"context"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryTracing(conf config.OtelConfiguration) *trace.TracerProvider {
	consoleExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	otlpExporter, err := getOtlTraceExporter(ctx, conf)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(consoleExporter),
		trace.WithBatcher(otlpExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OTEL_SERVICE_NAME),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp
}

func getOtlTraceExporter(ctx context.Context, conf config.OtelConfiguration) (*otlptrace.Exporter, error) {
	protocol := conf.OTEL_EXPORTER_OTLP_PROTOCOL
	if conf.OTEL_EXPORTER_OTLP_TRACES_PROTOCOL != "" {
		protocol = conf.OTEL_EXPORTER_OTLP_TRACES_PROTOCOL
	}

	if protocol == "http" {
		return otlptracehttp.New(ctx)
	}

	return otlptracegrpc.New(ctx)
}
