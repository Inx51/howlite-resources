package azureblob

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/resource"
	"github.com/inx51/howlite-resources/storage"
)

type Storage struct {
	containerClient *container.Client
	configuration   configuration.AzureBlobStorageConfiguration
}

func (azureBlobStorage *Storage) GetResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (*resource.Resource, error) {
	blobClient := azureBlobStorage.containerClient.NewBlobClient(resourceIdentifier.ToUniqueFilename())
	blobStream, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		blobStream.Body.Close()
		return nil, err
	}

	resource, err := resource.LoadResource(resourceIdentifier, blobStream.Body)

	return resource, err
}

func (azureBlobStorage *Storage) RemoveResource(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) error {
	blobClient := azureBlobStorage.containerClient.NewBlobClient(resourceIdentifier.ToUniqueFilename())
	_, err := blobClient.Delete(ctx, nil)

	return err
}

func (azureBlobStorage *Storage) ResourceExists(ctx context.Context, resourceIdentifier *resource.ResourceIdentifier) (bool, error) {
	blobClient := azureBlobStorage.containerClient.NewBlobClient(resourceIdentifier.ToUniqueFilename())
	_, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		if bloberror.HasCode(err, bloberror.BlobNotFound) {
			return false, nil
		}

		return false, err
	}

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
	blockBlobClient := azureBlobStorage.containerClient.NewBlockBlobClient(resource.Identifier.ToUniqueFilename())

	reader := azureBlobStorage.createResourceReader(resource)
	_, err := blockBlobClient.UploadStream(ctx, reader, &blockblob.UploadStreamOptions{
		BlockSize:   azureBlobStorage.configuration.BLOCK_SIZE,
		Concurrency: azureBlobStorage.configuration.UPLOAD_CONCURRENCY,
	})

	return err
}

func (azureBlobStorage *Storage) createResourceReader(resource *resource.Resource) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()
		err := resource.Write(pipeWriter)
		if err != nil {
			pipeWriter.CloseWithError(err)
		}
	}()
	return pipeReader
}
