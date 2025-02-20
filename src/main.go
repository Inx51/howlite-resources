package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"

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
	application.SetupOpenTelemetry()

	application.logger.Info("Starting application")

	application.SetupConfiguration()
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
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) SetupConfiguration() {

	app.logger.Info("Setting up configurations")

	godotenv.Overload(".env", ".env.local")

	config := config.Configuration{}
	env.Parse(&config)
	app.config = &config
}

func (app *Application) SetupOpenTelemetry() {

	serviceName, exists := os.LookupEnv("HOWLITE_SERVICE_NAME")
	if !exists {
		serviceName = "Howlite.Resrouces"
	}
	app.logger = telemetry.CreateOpenTelemetryLogger(serviceName)
}

func (app *Application) SetupStorage() {
	app.storage = filesystem.NewStorage(app.config.Path)
}

func (app *Application) SetupRepository() {
	app.repository = repository.NewRepository(&app.storage)
}

func (app *Application) SetupHandlers() {
	api.SetupHandlers(app.repository)
}

func (app *Application) Run() {
	api.Run(app.config.Host, app.config.Port)
}
