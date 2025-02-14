package main

import (
	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v11"
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage"
	"github.com/inx51/howlite/resources/storage/filesystem"
)

func main() {
	application := NewApplication()
	application.SetupConfiguration()
	application.SetupLogger()
	application.SetupStorage()
	application.SetupRepository()

	application.Run()
}

type Application struct {
	repository *repository.Repository
	storage    storage.Storage
	config     *config.Configuration
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

func (app *Application) SetupLogger() {
}

func (app *Application) SetupStorage() {
	app.storage = filesystem.NewStorage(app.config.Path)
}

func (app *Application) SetupRepository() {
	app.repository = repository.NewRepository(&app.storage)
}

func (app *Application) Run() {
	api.Run(app.repository)
}
