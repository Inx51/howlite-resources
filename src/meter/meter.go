package meter

import (
	"context"

	"github.com/inx51/howlite-resources/test"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter metric.Meter
var int64Counters map[string]metric.Int64Counter = make(map[string]metric.Int64Counter)

func SetupMeter() {
	if !test.IsTestRun() {
		meter = otel.Meter("howlite.resources")
	}
}

func ArithmeticInt64Counter(ctx context.Context, counterName string, change int64, options ...metric.AddOption) {
	if int64Counters[counterName] == nil {
		int64Counters[counterName], _ = meter.Int64Counter(counterName)
	}
	int64Counters[counterName].Add(ctx, change, options...)
}
