package main

import (
	"context"
	"time"

	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/http/handlers"
	"github.com/inx51/howlite-resources/http/server"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/storage/azureblob"
	"github.com/inx51/howlite-resources/storage/filesystem"
	"github.com/inx51/howlite-resources/storage/gcs"
	"github.com/inx51/howlite-resources/storage/s3"
)

type Container struct {
	storage  storage.Storage
	handlers *[]handlers.Handler
	server   *server.Server
}

func NewContainer() *Container {
	return &Container{}
}

func (container *Container) setupStorage(ctx context.Context, configuration configuration.StorageProvider) {
	storageProviderName := configuration.NAME
	switch storageProviderName {
	case "filesystem":
		container.storage = filesystem.NewStorage(&configuration.STORAGE_PROVIDER_FILESYSTEM)
	case "azureblob":
		container.storage = azureblob.NewStorage(&configuration.STORAGE_PROVIDER_AZBLOB)
	case "s3":
		container.storage = s3.NewStorage(&configuration.STORAGE_PROVIDER_S3)
	case "gcs":
		container.storage = gcs.NewStorage(&configuration.STORAGE_PROVIDER_GCS)
	default:
		panic("Unsupported storage provider: " + storageProviderName)
	}
	logger.Info(ctx, "Storage provider loaded", "provider", container.storage.GetName())
}

func (container *Container) setupHandlers() {
	container.handlers = &[]handlers.Handler{
		handlers.NewGetHandler(&container.storage),
		handlers.NewCreateHandler(&container.storage),
		handlers.NewReplaceHandler(&container.storage),
		handlers.NewRemoveHandler(&container.storage),
		handlers.NewExistsHandler(&container.storage),
		handlers.NewSysProbeHandler(),
	}
}

func (container *Container) setupHttpServer(configuration configuration.HttpServer) {

	readTimeout, err := time.ParseDuration(configuration.READ_TIMEOUT)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := time.ParseDuration(configuration.WRITE_TIMEOUT)
	if err != nil {
		panic(err)
	}
	idleTimeout, err := time.ParseDuration(configuration.IDLE_TIMEOUT)
	if err != nil {
		panic(err)
	}

	container.server = server.NewServer(
		configuration.HOST,
		configuration.PORT,
		container.handlers,
		readTimeout,
		writeTimeout,
		idleTimeout)
}
