package configuration

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Configuration struct {
	HTTP_SERVER      HttpServer
	STORAGE_PROVIDER StorageProvider
	OTEL             OtelConfiguration
	TRACING          Tracing
	EVENT_PUBLISHER  EventPublisher
}

type Tracing struct {
	LEVEL string `env:"HOWLITE_RESOURCE_TRACING_LEVEL" envDefault:"Info"`
}

type HttpServer struct {
	HOST          string `env:"HOWLITE_RESOURCE_HTTP_SERVER_HOST" envDefault:"localhost"`
	PORT          int    `env:"HOWLITE_RESOURCE_HTTP_SERVER_PORT" envDefault:"8080"`
	IDLE_TIMEOUT  string `env:"HOWLITE_RESOURCE_HTTP_SERVER_IDLE_TIMEOUT" envDefault:"30s"`
	READ_TIMEOUT  string `env:"HOWLITE_RESOURCE_HTTP_SERVER_READ_TIMEOUT" envDefault:"30s"`
	WRITE_TIMEOUT string `env:"HOWLITE_RESOURCE_HTTP_SERVER_WRITE_TIMEOUT" envDefault:"30s"`
}

type OtelConfiguration struct {
	OTEL_SERVICE_NAME                   string `env:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_PROTOCOL         string `env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	OTEL_EXPORTER_OTLP_TRACES_PROTOCOL  string `env:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"`
	OTEL_EXPORTER_OTLP_METRICS_PROTOCOL string `env:"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"`
	OTEL_EXPORTER_OTLP_LOGS_PROTOCOL    string `env:"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"`
}

type StorageProvider struct {
	NAME                        string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME" envDefault:"filesystem"`
	STORAGE_PROVIDER_FILESYSTEM FilesystemConfiguration
	STORAGE_PROVIDER_S3         S3Configuration
	STORAGE_PROVIDER_AZBLOB     AzureBlobStorageConfiguration
}

type FilesystemConfiguration struct {
	PATH string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_FILESYSTEM_PATH" envDefault:"./tmp/howlite"`
}

type S3Configuration struct {
	BUCKET               string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_BUCKET"`
	ACCESS_KEY           string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ACCESS_KEY"`
	SECRET_KEY           string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_SECRET_KEY"`
	ENDPOINT             string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ENDPOINT"`
	REGION               string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_REGION"`
	PART_UPLOAD_SIZE     int64  `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_PART_UPLOAD_SIZE" envDefault:"5242880"`
	UPLOAD_CONCURRENCY   int    `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_UPLOAD_CONCURRENCY" envDefault:"5"`
	DOWNLOAD_CONCURRENCY int    `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_DOWNLOAD_CONCURRENCY" envDefault:"5"`
}

type AzureBlobStorageConfiguration struct {
	CONNECTION_STRING  string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONNECTION_STRING"`
	CONTAINER_NAME     string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONTAINER_NAME"`
	BLOCK_SIZE         int64  `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_BLOCK_SIZE" envDefault:"8388608"`
	UPLOAD_CONCURRENCY int    `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_UPLOAD_CONCURRENCY" envDefault:"5"`
}

type EventPublisher struct {
	OUTBOX_SQLITE_PATH   string `env:"HOWLITE_RESOURCE_EVENT_PUBLISHER_OUTBOX_SQLITE_PATH"`
	ZEROMQ_CONFIGURATION ZeroMqConfiguration
}

// SERVER_CERT_PATH must point at a CZMQ secret cert file (the "*_secret"
// file produced by zcert_save / goczmq-certgen), which holds both the
// publisher's public and secret key.
//
// ALLOWED_CLIENTS_PATH is optional and points at a directory of public
// cert files (a CZMQ certstore) for the subscribers allowed to connect. If
// left empty, any client with a CURVE keypair is accepted (CURVE_ALLOW_ANY) -
// connections are still encrypted, but not restricted to known peers.
type ZeroMqConfiguration struct {
	ENDPOINT string `env:"HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_ENDPOINT"`
	CURVE    ZeroMqCurveConfiguration
}

type ZeroMqCurveConfiguration struct {
	SERVER_CERT_PATH     string `env:"HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_CURVE_SERVER_CERT_PATH"`
	ALLOWED_CLIENTS_PATH string `env:"HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_CURVE_ALLOWED_CLIENTS_PATH"`
}

//TODO: We should validate the configuration values so we can throw any unexpected configuration errors on startup..

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func ConfigureEnvFiles() {
	godotenv.Overload(".env", ".local.env")
}

func ConfigureEnvironmentVariables(configuration *Configuration) {
	env.Parse(configuration)
}
