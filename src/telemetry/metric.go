package telemetry

import (
	"context"

	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/meter"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
)

var meterProvider *metric.MeterProvider

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	otlpMetricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		logger.Warn(ctx, "Failed to create metrics exporter for OpenTelemetry, skipping OpenTelemetry metrics", "error", err)
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(otlpMetricReader),
	)
	return provider, nil
}

func SetupMetric(ctx context.Context) {
	var err error
	meterProvider, err = newMeterProvider(ctx)
	if err != nil {
		return
	}

	otel.SetMeterProvider(meterProvider)
	meter.SetupMeter(meterProvider != nil)
}

func ShutdownMetrics(ctx context.Context) {
	if err := meterProvider.Shutdown(ctx); err != nil {
		logger.Error(ctx, "OpenTelemetry failed to shutdown meter provider", "error", err)
		return
	}
}
