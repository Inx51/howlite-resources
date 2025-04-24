package telemetry

import (
	"context"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryMeter(conf config.OtelConfiguration) *metric.MeterProvider {
	ctx := context.Background()

	otlpExporter, err := getOtlpMeterExporter(ctx, conf)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(otlpExporter)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OTEL_SERVICE_NAME),
		)))

	return meterProvider
}

func getOtlpMeterExporter(ctx context.Context, conf config.OtelConfiguration) (metric.Exporter, error) {
	protocol := conf.OTEL_EXPORTER_OTLP_PROTOCOL
	if conf.OTEL_EXPORTER_OTLP_METRICS_PROTOCOL != "" {
		protocol = conf.OTEL_EXPORTER_OTLP_METRICS_PROTOCOL
	}

	if strings.HasPrefix(protocol, "http") {
		return otlpmetricgrpc.New(ctx)
	}

	return otlpmetrichttp.New(ctx)
}
