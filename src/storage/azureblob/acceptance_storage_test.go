package azureblob

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/http/handlers"
	httpserver "github.com/inx51/howlite-resources/http/server"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

const testContainer = "test-container"

func newTestServer(t *testing.T) (*httptest.Server, *http.Client) {
	t.Helper()
	ctx := context.Background()

	ctr, err := tcazurite.Run(
		ctx,
		"mcr.microsoft.com/azure-storage/azurite:3.28.0",
		tcazurite.WithEnabledServices(tcazurite.BlobService),
		testcontainers.WithCmdArgs("--skipApiVersionCheck"),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, ctr)

	blobURL, err := ctr.BlobServiceURL(ctx)
	require.NoError(t, err)

	connString := fmt.Sprintf(
		"DefaultEndpointsProtocol=http;AccountName=%s;AccountKey=%s;BlobEndpoint=%s/%s;",
		tcazurite.AccountName,
		tcazurite.AccountKey,
		blobURL,
		tcazurite.AccountName,
	)

	azblobClient, err := azblob.NewClientFromConnectionString(connString, nil)
	require.NoError(t, err)
	_, err = azblobClient.CreateContainer(ctx, testContainer, nil)
	require.NoError(t, err)

	storageConfig := &configuration.AzureBlobStorageConfiguration{
		CONNECTION_STRING:  connString,
		CONTAINER_NAME:     testContainer,
		BLOCK_SIZE:         8388608,
		UPLOAD_CONCURRENCY: 5,
	}

	store := NewStorage(storageConfig)
	hs := &[]handlers.Handler{
		handlers.NewGetHandler(&store),
		handlers.NewCreateHandler(&store),
		handlers.NewReplaceHandler(&store),
		handlers.NewRemoveHandler(&store),
		handlers.NewExistsHandler(&store),
	}

	ts := httptest.NewServer(httpserver.NewServeMux(hs))
	t.Cleanup(ts.Close)
	return ts, ts.Client()
}

func TestAcceptance_GetResource_ReturnsResource(t *testing.T) {
	ts, client := newTestServer(t)

	postResp, err := client.Post(ts.URL+"/my/resource.txt", "text/plain", strings.NewReader("hello world"))
	require.NoError(t, err)
	postResp.Body.Close()

	getResp, err := client.Get(ts.URL + "/my/resource.txt")
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)
	got, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(got))
}

func TestAcceptance_GetResource_ReturnsNotFoundWhenResourceDoesNotExist(t *testing.T) {
	ts, client := newTestServer(t)

	resp, err := client.Get(ts.URL + "/does/not/exist.txt")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAcceptance_ExistsResource_ReturnsNoContentWhenResourceExists(t *testing.T) {
	ts, client := newTestServer(t)

	postResp, err := client.Post(ts.URL+"/my/resource.txt", "text/plain", strings.NewReader("hello world"))
	require.NoError(t, err)
	postResp.Body.Close()

	req, _ := http.NewRequest(http.MethodHead, ts.URL+"/my/resource.txt", nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAcceptance_ExistsResource_ReturnsNotFoundWhenResourceDoesNotExist(t *testing.T) {
	ts, client := newTestServer(t)

	req, _ := http.NewRequest(http.MethodHead, ts.URL+"/my/resource.txt", nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAcceptance_RemoveResource_RemovesResource(t *testing.T) {
	ts, client := newTestServer(t)

	postResp, err := client.Post(ts.URL+"/my/resource.txt", "text/plain", strings.NewReader("hello world"))
	require.NoError(t, err)
	postResp.Body.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/my/resource.txt", nil)
	delResp, err := client.Do(req)
	require.NoError(t, err)
	delResp.Body.Close()
	require.Equal(t, http.StatusNoContent, delResp.StatusCode)

	getResp, err := client.Get(ts.URL + "/my/resource.txt")
	require.NoError(t, err)
	getResp.Body.Close()
	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestAcceptance_RemoveResource_ReturnsNotFoundWhenResourceDoesNotExist(t *testing.T) {
	ts, client := newTestServer(t)

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/my/resource.txt", nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAcceptance_ReplaceResource_CreatesResourceWhenResourceDoesNotExist(t *testing.T) {
	ts, client := newTestServer(t)

	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/my/resource.txt", strings.NewReader("version one"))
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotEmpty(t, resp.Header.Get("Location"))

	getResp, err := client.Get(ts.URL + "/my/resource.txt")
	require.NoError(t, err)
	body, err := io.ReadAll(getResp.Body)
	getResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "version one", string(body))
}

func TestAcceptance_ReplaceResource_ReplacesExistingResource(t *testing.T) {
	ts, client := newTestServer(t)

	postResp, err := client.Post(ts.URL+"/my/resource.txt", "text/plain", strings.NewReader("version one"))
	require.NoError(t, err)
	postResp.Body.Close()

	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/my/resource.txt", strings.NewReader("version two"))
	req.Header.Set("Content-Type", "text/plain")
	replaceResp, err := client.Do(req)
	require.NoError(t, err)
	replaceResp.Body.Close()
	require.Equal(t, http.StatusNoContent, replaceResp.StatusCode)
	require.NotEmpty(t, replaceResp.Header.Get("Location"))

	getResp, err := client.Get(ts.URL + "/my/resource.txt")
	require.NoError(t, err)
	body, err := io.ReadAll(getResp.Body)
	getResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "version two", string(body))
}

func TestAcceptance_CreateResource(t *testing.T) {
	ts, client := newTestServer(t)

	body := strings.NewReader("hello world")
	resp, err := client.Post(ts.URL+"/my/resource.txt", "text/plain", body)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotEmpty(t, resp.Header.Get("Location"))

	getResp, err := client.Get(ts.URL + "/my/resource.txt")
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)

	got, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(got))
}
