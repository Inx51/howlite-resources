package s3

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/inx51/howlite-resources/configuration"
	"github.com/inx51/howlite-resources/http/handlers"
	httpserver "github.com/inx51/howlite-resources/http/server"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

const testBucket = "test-bucket"

func TestMain(m *testing.M) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) (*httptest.Server, *http.Client) {
	t.Helper()
	ctx := context.Background()

	ctr, err := tcminio.Run(ctx, "minio/minio:RELEASE.2024-01-16T16-07-38Z")
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, ctr)

	endpoint, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	const region = "us-east-1"
	endpointURL := "http://" + endpoint
	accessKey, secretKey := ctr.Username, ctr.Password

	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithBaseEndpoint(endpointURL),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	require.NoError(t, err)

	s3Client := awss3.NewFromConfig(cfg, func(o *awss3.Options) {
		o.UsePathStyle = true
	})
	_, err = s3Client.CreateBucket(ctx, &awss3.CreateBucketInput{
		Bucket: aws.String(testBucket),
	})
	require.NoError(t, err)

	storageConfig := &configuration.S3Configuration{
		BUCKET:               testBucket,
		ACCESS_KEY:           accessKey,
		SECRET_KEY:           secretKey,
		ENDPOINT:             endpointURL,
		REGION:               region,
		PART_UPLOAD_SIZE:     5242880,
		UPLOAD_CONCURRENCY:   5,
		DOWNLOAD_CONCURRENCY: 5,
	}

	store := NewStorage(ctx, storageConfig)
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
