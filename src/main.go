package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/caarlos0/env/v11"
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/storage/filesystem"
	"github.com/inx51/howlite/resources/telemetry"
)

func main() {
	application := NewApplication()
	application.SetupConfiguration()
	application.SetupOpenTelemetry()
	application.SetupStorage()
	application.SetupRepository()
	application.SetupHandlers()

	application.Run()
}

type Application struct {
	repository *repository.Repository
	storage    storage.Storage
	config     *config.Configuration
	logger     *slog.Logger
	tracer     *trace.TracerProvider
	meter      *metric.MeterProvider
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
	app.tracer = telemetry.CreateOpenTelemetryTracer(app.config.OTEL)
	app.meter = telemetry.CreateOpenTelemetryMeter(app.config.OTEL)
}

func (app *Application) SetupStorage() {
	app.logger.Debug("Trying to setup storage")
	app.storage = filesystem.NewStorage(app.config.PATH, app.logger)
	app.logger.Info("Setup storage provider", "provider", app.storage.GetName())
}

func (app *Application) SetupRepository() {
	app.repository = repository.NewRepository(&app.storage, app.logger)
}

func (app *Application) SetupHandlers() {
	api.SetupHandlers(
		app.repository,
		app.logger,
		app.meter)
}

func (app *Application) Run() {
	api.Run(app.config.HOST, app.config.PORT, app.logger)
}
