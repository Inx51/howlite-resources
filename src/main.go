package main

import (
	"context"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/telemetry"
	"github.com/joho/godotenv"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	application := NewApplication()
	application.SetupConfiguration()
	application.SetupOpenTelemetry()
	defer (*application.span).End()

	application.SetupStorageContext(ctx)
	application.SetupRepository()
	application.SetupHandlers()

	application.RunContext(ctx)
}

type Application struct {
	repository *repository.Repository
	storage    *storage.Storage
	config     *config.Configuration
	logger     *slog.Logger
	meter      *otelmetric.Meter
	span       *oteltrace.Span
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) SetupConfiguration() {

	godotenv.Overload(".env", ".env.local")

	config := config.Configuration{}
	env.Parse(&config)
	app.config = &config
}

func (app *Application) SetupOpenTelemetry() {
	app.logger = telemetry.CreateOpenTelemetryLogger(app.config.OTEL)
	app.span = telemetry.CreateOpenTelemetryTracer(app.config.OTEL)
	app.meter = telemetry.CreateOpenTelemetryMeter(app.config.OTEL)
}

func (app *Application) SetupStorageContext(ctx context.Context) {
	app.logger.DebugContext(ctx, "Trying to setup storage")
	storage, err := storage.Create(app.logger, app.config.STORAGE_PROVIDER)
	if err != nil {
		app.logger.ErrorContext(ctx, "Failed to setup storage", "error", err)
	}
	app.storage = &storage
	// filesystem.NewStorage(app.config.PATH, app.logger)
	app.logger.InfoContext(ctx, "Setup storage provider", "provider", (*app.storage).GetName())
}

func (app *Application) SetupRepository() {
	app.repository = repository.NewRepository(app.storage, app.logger)
}

func (app *Application) SetupHandlers() {
	api.SetupHandlers(
		app.repository,
		app.logger,
		app.meter)
}

func (app *Application) RunContext(ctx context.Context) {
	api.RunContext(ctx, app.config.HOST, app.config.PORT, app.logger)
}
