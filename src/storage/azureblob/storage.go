package azureblob

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
)

type Storage struct {
	containerClient *container.Client
	configuration   configuration.AzureBlobStorageConfiguration
}

func (azureBlobStorage *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	blobName := resourceIdentifier.ToUniqueFilename()
	logger.Debug(ctx, "trying to download blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
	blobClient := azureBlobStorage.containerClient.NewBlobClient(blobName)

	blobClientCtx, span := tracer.StartDebugSpan(ctx, "azure.blob.download")
	tracer.SetDebugAttributes(blobClientCtx, span,
		attribute.String("blob.name", blobName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	blobStream, err := blobClient.DownloadStream(blobClientCtx, nil)
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			logger.Debug(ctx, "blob not found", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
			return nil, err
		}
		logger.Error(ctx, "failed to download blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName, "error", err)
		return nil, err
	}

	resource, err := resource.LoadResource(resourceIdentifier, blobStream.Body)
	if err != nil {
		blobStream.Body.Close()
		logger.Error(ctx, "failed to load resource from blob stream", "resource.identifier", resourceIdentifier.Identifier(), "error", err)
		return nil, err
	}
	logger.Debug(ctx, "successfully downloaded blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
	return resource, err
}

func (azureBlobStorage *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	blobName := resourceIdentifier.ToUniqueFilename()
	blobClient := azureBlobStorage.containerClient.NewBlobClient(blobName)
	logger.Debug(ctx, "trying to delete blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)

	blobClientCtx, span := tracer.StartDebugSpan(ctx, "azure.blob.delete")
	tracer.SetDebugAttributes(blobClientCtx, span,
		attribute.String("blob.name", blobName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	_, err := blobClient.Delete(blobClientCtx, nil)
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		logger.Error(ctx, "failed to delete blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName, "error", err)
		return err
	}
	logger.Debug(ctx, "successfully deleted blob", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
	return nil
}

func (azureBlobStorage *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	blobName := resourceIdentifier.ToUniqueFilename()
	blobClient := azureBlobStorage.containerClient.NewBlobClient(blobName)
	logger.Debug(ctx, "checking if blob exists", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)

	blobClientCtx, span := tracer.StartDebugSpan(ctx, "azure.blob.get_properties")
	tracer.SetDebugAttributes(blobClientCtx, span,
		attribute.String("blob.name", blobName),
		attribute.String("resource.identifier", resourceIdentifier.Identifier()),
	)
	_, err := blobClient.GetProperties(blobClientCtx, nil)
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			logger.Debug(ctx, "blob not found", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
			return false, nil
		}

		logger.Error(ctx, "failed to get blob properties", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName, "error", err)
		return false, err
	}

	logger.Debug(ctx, "blob exists", "resource.identifier", resourceIdentifier.Identifier(), "blob.name", blobName)
	return true, nil
}

func NewStorage(configuration *configuration.AzureBlobStorageConfiguration) storage.Storage {
	client, err := azblob.NewClientFromConnectionString(configuration.CONNECTION_STRING, nil)
	if err != nil {
		panic(err)
	}

	return &Storage{
		configuration:   *configuration,
		containerClient: client.ServiceClient().NewContainerClient(configuration.CONTAINER_NAME),
	}
}

func (azureBlobStorage *Storage) GetName() string {
	return "azureblob"
}

func (azureBlobStorage *Storage) SaveResource(ctx context.Context, resource *resource.Resource) error {
	blobName := resource.Identifier.ToUniqueFilename()
	blockBlobClient := azureBlobStorage.containerClient.NewBlockBlobClient(blobName)

	reader := azureBlobStorage.createResourceReader(ctx, resource)
	logger.Debug(ctx, "trying to upload blob", "resource.identifier", resource.Identifier.Identifier(), "blob.name", blobName)

	blobClientCtx, span := tracer.StartDebugSpan(ctx, "azure.blob.upload")
	tracer.SetDebugAttributes(blobClientCtx, span,
		attribute.String("blob.name", blobName),
		attribute.String("resource.identifier", resource.Identifier.Identifier()),
		attribute.Int64("azure.blob.block_size", azureBlobStorage.configuration.BLOCK_SIZE),
		attribute.Int("azure.blob.concurrency", azureBlobStorage.configuration.UPLOAD_CONCURRENCY),
	)
	_, err := blockBlobClient.UploadStream(blobClientCtx, reader, &blockblob.UploadStreamOptions{
		BlockSize:   azureBlobStorage.configuration.BLOCK_SIZE,
		Concurrency: azureBlobStorage.configuration.UPLOAD_CONCURRENCY,
	})
	if err != nil {
		span.RecordError(err)
	}
	tracer.SafeEndSpan(span)

	if err != nil {
		logger.Error(ctx, "failed to upload blob", "resource.identifier", resource.Identifier.Identifier(), "blob.name", blobName, "error", err)
		return err
	}

	logger.Debug(ctx, "successfully uploaded blob", "resource.identifier", resource.Identifier.Identifier(), "blob.name", blobName)
	return nil
}

func (azureBlobStorage *Storage) createResourceReader(ctx context.Context, resource *resource.Resource) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()
		err := resource.Write(pipeWriter)
		if err != nil {
			logger.Error(ctx, "failed to write resource to pipe", "error", err)
			pipeWriter.CloseWithError(err)
		}
	}()
	return pipeReader
}
