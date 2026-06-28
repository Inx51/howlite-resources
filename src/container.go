package main

import (
	"context"
	"time"

	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/event"
	"github.com/inx51/howlite-resources/http/handlers"
	"github.com/inx51/howlite-resources/http/server"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/storage/azureblob"
	"github.com/inx51/howlite-resources/storage/filesystem"
	"github.com/inx51/howlite-resources/storage/s3"
)

type Container struct {
	storage      storage.Storage
	handlers     *[]handlers.Handler
	server       *server.Server
	bus          *event.Bus
	outboxWorker *event.OutboxWorker
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
		container.storage = s3.NewStorage(ctx, &configuration.STORAGE_PROVIDER_S3)
	default:
		panic("Unsupported storage provider: " + storageProviderName)
	}
	logger.Info(ctx, "Storage provider loaded", "provider", container.storage.GetName())
}

func (container *Container) setupHandlers() {
	container.handlers = &[]handlers.Handler{
		handlers.NewGetHandler(&container.storage),
		handlers.NewCreateHandler(&container.storage, container.bus),
		handlers.NewReplaceHandler(&container.storage, container.bus),
		handlers.NewRemoveHandler(&container.storage, container.bus),
		handlers.NewExistsHandler(&container.storage),
		handlers.NewSysProbeHandler(),
	}
}

func (container *Container) setupEventPublisher(ctx context.Context, configuration configuration.EventPublisher) {

	var publisherPtr *event.Publisher
	var outboxPtr *event.Outbox

	if configuration.EVENT_PUBLISHER_ENDPOINT != "" {
		publisher := event.NewPublisher(ctx, configuration.EVENT_PUBLISHER_ENDPOINT)
		if !publisher.IsAvailable() {
			logger.Error(ctx, "Event publisher configured but unavailable in this build")
			return
		}
		publisherPtr = &publisher
	} else {
		logger.Info(ctx, "No event publisher endpoint specified, events will not be published")
		return
	}

	if configuration.OUTBOX_SQLITE_PATH != "" {
		outbox := event.NewOutbox(ctx, configuration.OUTBOX_SQLITE_PATH)
		outboxPtr = &outbox

		outboxWorker := event.NewOutboxWorker(ctx, outboxPtr, publisherPtr)
		container.outboxWorker = &outboxWorker
	} else {
		logger.Info(ctx, "No outbox path specified, published events will not be persisted")
	}

	container.bus = event.NewBus(publisherPtr, outboxPtr) // One or both can be nil
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
