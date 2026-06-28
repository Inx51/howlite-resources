
<p>
	<img src="https://img.shields.io/badge/status-in%20beta-orange" alt="status"/>
</p>

> [!IMPORTANT]  
> All code and documentation is subject to change and is currently in beta!

# Howlite Resources

> Simple, pluggable HTTP resource store for fast prototyping and storage of any type of data.

Howlite Resources lets you store, fetch, update, and delete any kind of resource (files, blobs, data) over HTTP. It's ideal for rapid prototyping, microservices, or as a backend for storing any type of data.

---

## ✨ Features

- **RESTful API:** POST, GET, PUT, DELETE, HEAD for resources
- **Pluggable storage:** Filesystem, S3, Azure Blob Storage
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

#### S3

Store resources in an S3-compatible object storage (e.g., AWS S3, MinIO).

Uploads and downloads are performed in parallel parts/chunks. Increasing concurrency can improve throughput for large resources by transferring multiple parts simultaneously, but each concurrent part consumes memory and a goroutine for its duration. The part size determines how the resource is split — smaller parts mean more parallel operations (up to the concurrency limit); larger parts reduce overhead but require more memory per part.

> **Rule of thumb:** The default values work well for most cases. Raise concurrency only if you have large resources, high network bandwidth, and can afford the added memory usage per request. Each upload/download request uses up to `concurrency × part_size` bytes of memory.

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME | s3 | Selects the S3 provider |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_BUCKET |  | S3 bucket name |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ACCESS_KEY |  | S3 access key |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_SECRET_KEY |  | S3 secret key |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ENDPOINT |  | S3 endpoint URL (leave empty for AWS) |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_REGION |  | S3 region |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_PART_UPLOAD_SIZE | 5242880 | Size of each part in a multipart transfer (bytes). Affects both upload and download chunking. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_UPLOAD_CONCURRENCY | 5 | Number of parts uploaded in parallel per PUT/POST request. Higher values increase upload speed for large resources at the cost of memory and CPU. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_DOWNLOAD_CONCURRENCY | 5 | Number of parts downloaded in parallel per GET request. Higher values increase download speed for large resources at the cost of memory and CPU. |

#### Azure Blob Storage

Store resources in Azure Blob Storage. Uploads are split into blocks that are sent in parallel and then committed as a single blob.

The block size determines how the resource is divided for upload. Larger blocks mean fewer network round-trips but higher memory usage per upload. Concurrency controls how many blocks are in-flight at once — increasing it can improve upload speed for large blobs, at the cost of more memory (up to `concurrency × block_size` bytes per upload request).

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME | azureblob | Selects the Azure Blob Storage provider |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONNECTION_STRING |  | Azure Storage connection string |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONTAINER_NAME |  | Blob container name |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_BLOCK_SIZE | 8388608 | Size of each block in a block blob upload (bytes, default 8 MiB). Larger values reduce round-trips but increase memory usage per upload. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_UPLOAD_CONCURRENCY | 5 | Number of blocks uploaded in parallel per PUT/POST request. Higher values increase upload speed for large blobs at the cost of memory and CPU. |

### Event Publisher

Howlite Resources can publish events when resources are created, replaced, or removed. Events are delivered reliably via a SQLite-backed outbox.

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_EVENT_PUBLISHER_ENDPOINT |  | ZeroMQ endpoint to publish events to |
| HOWLITE_RESOURCE_EVENT_OUTBOX_SQLITE_PATH |  | Path to the SQLite outbox database file |

### Telemetry

| Variable | Default | Description |
|---|---|---|
| HOWLITE_RESOURCE_TRACING_LEVEL | Info | Tracing level |

Support for standard OTEL environment variables.
Read more [here](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/)

---

## 📄 License

[MIT](LICENSE)
