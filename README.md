
<p>
	<img src="https://img.shields.io/badge/status-in%20development-orange" alt="status"/>
</p>

# Howlite Resources

> Simple, pluggable HTTP resource store for fast prototyping and storage of any type of data.

Howlite Resources lets you store, fetch, update, and delete any kind of resource (files, blobs, data) over HTTP. It's ideal for rapid prototyping, microservices, or as a backend for storing any type of data.

---

## ✨ Features

- **RESTful API:** POST, GET, PUT, DELETE, HEAD for resources
- **Pluggable storage:** Filesystem, S3, and more
- **OpenTelemetry:** Metrics & tracing built-in
- **Easy config:** Environment variables or .env

---

## 🚀 Quick Start

```
!! DOCKER IMAGE WILL BE AVAILABLE LATER FOR USE !!
!! BELOW INSTRUCTIONS MIGHT CHANGE !!
```

```sh
git clone https://github.com/Inx51/howlite-resources.git
cd src
go build -o howlite-resources ./src
./howlite-resources
```

---

## 🛣️ API Overview

| Method | Path | Description |
|--------|------|-------------|
| POST   | /your/resource/path    | Create resource |
| GET    | /your/resource/path    | Get resource    |
| PUT    | /your/resource/path    | Replace/create  |
| DELETE | /your/resource/path    | Remove resource |
| HEAD   | /your/resource/path    | Resource exists |

---

## ⚙️ Configuration

Set via environment variables.

### HTTP Server

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_HTTP_SERVER_HOST | localhost | HTTP server host |
| HOWLITE_RESOURCE_HTTP_SERVER_PORT | 8080 | HTTP server port |
| HOWLITE_RESOURCE_HTTP_SERVER_IDLE_TIMEOUT | 30s | Idle timeout of request/response |
| HOWLITE_RESOURCE_HTTP_SERVER_READ_TIMEOUT | 30s | Read timeout of request/response |
| HOWLITE_RESOURCE_HTTP_SERVER_WRITE_TIMEOUT | 30s | Write timeout of request/response |


### Storage Providers

Howlite Resources supports multiple storage backends. Select the provider with `HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME`.

#### Filesystem

The default provider. Stores resources as files on the local disk.

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME | filesystem | Selects the filesystem provider |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_FILESYSTEM_PATH | ./tmp/howlite | Directory for storing files |

#### S3 (Work in progress)

Store resources in an S3-compatible object storage (e.g., AWS S3, MinIO).

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME | s3 | Selects the S3 provider |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_BUCKET |  | S3 bucket name |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_PREFIX |  | Prefix for S3 object keys |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ACCESS_KEY |  | S3 access key |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_SECRET_KEY |  | S3 secret key |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ENDPOINT |  | S3 endpoint URL |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_REGION |  | S3 region |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_UPLOAD_STRATEGY | singlepart | Upload strategy: singlepart or multipart |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_MULTIPART_PART_UPLOAD_SIZE | 5242880 | Multipart upload part size (bytes) |

<!-- Add more providers here as they are implemented -->

### Telemetry

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_TRACING_LEVEL | Debug, Info | Tracing level |

Support for standard OTEL environemnt variables.
Read more [here](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/)

---

## 📄 License

[MIT](LICENSE)
