package telemetry

import (
	"context"

	"github.com/inx51/howlite/resources/config"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func CreateOpenTelemetryMeter(conf config.OtelConfiguration) *otelmetric.Meter {
	ctx := context.Background()

	otlpMetricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(otlpMetricReader),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OTEL_SERVICE_NAME),
		)))

	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter("howlite.resources")

	return &meter
}
