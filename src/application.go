package main

import (
	"context"

	"github.com/inx51/howlite-resources/configuration"
)

type Application struct {
	container     *Container
	configuration *configuration.Configuration
}

func NewApplication() *Application {
	return &Application{
		container:     NewContainer(),
		configuration: configuration.NewConfiguration(),
	}
}

func (app *Application) ConfigureConfigurations(ctx context.Context) *configuration.Configuration {
	configurations := configuration.NewConfiguration()

	configuration.ConfigureEnvFiles()
	configuration.ConfigureEnvironmentVariables(configurations)

	app.configuration = configurations
	return configurations
}

func (app *Application) ConfigureContainer(ctx context.Context) {
	container := NewContainer()
	container.setupStorage(ctx, app.configuration.STORAGE_PROVIDER)
	container.setupHandlers()
	container.setupHttpServer(app.configuration.HTTP_SERVER)

	app.container = container
}

func (app *Application) Run(ctx context.Context) {
	app.container.server.Start(ctx)
}

func (app *Application) Shutdown(ctx context.Context) {
	app.container.server.Shutdown(ctx)
}
