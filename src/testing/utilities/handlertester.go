package utilities

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/testing/utilities/logging"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func CreateHandlerParameters(tester *tester.Tester, storage storage.Storage) (context.Context, http.ResponseWriter, *http.Request, *repository.Repository, *slog.Logger, *otelmetric.Meter) {
	resp, req := tester.Build()
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	repo := repository.NewRepository(&storage, logger)
	meter := metric.NewMeterProvider().Meter("test")

	return context.Background(), resp, req, repo, logger, &meter
}
