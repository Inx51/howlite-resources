package utilities

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	"go.opentelemetry.io/otel/sdk/metric"
)

func CreateHandlerParameters(tester *tester.Tester, storage storage.Storage) (context.Context, http.ResponseWriter, *http.Request, *repository.Repository, *slog.Logger, *metric.MeterProvider) {
	resp, req := tester.Build()
	logger := slog.New(slog.NewTextHandler(&TestingLogWriter{}, nil))
	repo := repository.NewRepository(&storage, logger)
	meter := metric.NewMeterProvider()

	return context.Background(), resp, req, repo, logger, meter
}
