
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
- **Event publishing:** Optional ZeroMQ events on resource changes, with CURVE encryption and a SQLite-backed outbox for reliable delivery
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

All variables are optional and have working defaults.

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_HTTP_SERVER_HOST | No | localhost | HTTP server host |
| HOWLITE_RESOURCE_HTTP_SERVER_PORT | No | 8080 | HTTP server port |
| HOWLITE_RESOURCE_HTTP_SERVER_IDLE_TIMEOUT | No | 30s | Idle timeout of request/response |
| HOWLITE_RESOURCE_HTTP_SERVER_READ_TIMEOUT | No | 30s | Read timeout of request/response |
| HOWLITE_RESOURCE_HTTP_SERVER_WRITE_TIMEOUT | No | 30s | Write timeout of request/response |


### Storage Providers

Howlite Resources supports multiple storage backends, selected by a single variable:

| Variable | Required | Default | Valid values | Description |
|---|---|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_NAME | No | filesystem | `filesystem`, `s3`, `azureblob` | Selects the storage provider |

Each provider has its own additional configuration below.

#### Filesystem

The default provider. Stores resources as files on the local disk.

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_FILESYSTEM_PATH | No | ./tmp/howlite | Directory for storing files |

#### S3

Store resources in an S3-compatible object storage (e.g., AWS S3, MinIO).

Uploads and downloads are performed in parallel parts/chunks. Increasing concurrency can improve throughput for large resources by transferring multiple parts simultaneously, but each concurrent part consumes memory and a goroutine for its duration. The part size determines how the resource is split — smaller parts mean more parallel operations (up to the concurrency limit); larger parts reduce overhead but require more memory per part.

> **Rule of thumb:** The default values work well for most cases. Raise concurrency only if you have large resources, high network bandwidth, and can afford the added memory usage per request. Each upload/download request uses up to `concurrency × part_size` bytes of memory.

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_BUCKET | Yes |  | S3 bucket name |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ACCESS_KEY | No |  | S3 access key. Leave empty (together with `SECRET_KEY`) to fall back to the standard AWS credential chain (env vars, shared config/profile, IAM role) |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_SECRET_KEY | No |  | S3 secret key. Leave empty (together with `ACCESS_KEY`) to fall back to the standard AWS credential chain |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ENDPOINT | No |  | S3 endpoint URL. Leave empty for AWS; set for S3-compatible services (e.g. MinIO) |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_REGION | No |  | S3 region. Leave empty to fall back to the standard AWS region resolution (env vars, shared config/profile) |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_PART_UPLOAD_SIZE | No | 5242880 | Size of each part in a multipart transfer (bytes). Affects both upload and download chunking. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_UPLOAD_CONCURRENCY | No | 5 | Number of parts uploaded in parallel per PUT/POST request. Higher values increase upload speed for large resources at the cost of memory and CPU. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_DOWNLOAD_CONCURRENCY | No | 5 | Number of parts downloaded in parallel per GET request. Higher values increase download speed for large resources at the cost of memory and CPU. |

#### Azure Blob Storage

Store resources in Azure Blob Storage. Uploads are split into blocks that are sent in parallel and then committed as a single blob.

The block size determines how the resource is divided for upload. Larger blocks mean fewer network round-trips but higher memory usage per upload. Concurrency controls how many blocks are in-flight at once — increasing it can improve upload speed for large blobs, at the cost of more memory (up to `concurrency × block_size` bytes per upload request).

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONNECTION_STRING | Yes |  | Azure Storage connection string |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_CONTAINER_NAME | Yes |  | Blob container name |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_BLOCK_SIZE | No | 8388608 | Size of each block in a block blob upload (bytes, default 8 MiB). Larger values reduce round-trips but increase memory usage per upload. |
| HOWLITE_RESOURCE_STORAGE_PROVIDER_AZUREBLOB_UPLOAD_CONCURRENCY | No | 5 | Number of blocks uploaded in parallel per PUT/POST request. Higher values increase upload speed for large blobs at the cost of memory and CPU. |

### Event Publisher

Howlite Resources can publish events over ZeroMQ when resources are created, replaced, or removed.

Everything below is optional — by default no `HOWLITE_RESOURCE_EVENT_PUBLISHER_*` variables are set, and event publishing is fully disabled. Each feature turns on as soon as its one "trigger" variable is set:

- **Event publishing** turns on once `ZEROMQ_ENDPOINT` is set.
- **Outbox persistence** turns on once `OUTBOX_SQLITE_PATH` is set — but only has an effect if event publishing is also enabled.
- **CURVE encryption** turns on once `ZEROMQ_CURVE_SERVER_CERT_PATH` is set — but only has an effect if event publishing is also enabled.

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_ENDPOINT | No — leave empty to disable event publishing |  | ZeroMQ endpoint to publish events to. Setting this is what turns event publishing on. |
| HOWLITE_RESOURCE_EVENT_PUBLISHER_OUTBOX_SQLITE_PATH | No |  | Path to the SQLite outbox database file. Leave empty to publish events directly with no persistence. Ignored if `ZEROMQ_ENDPOINT` is not set. |

#### CURVE (transport security)

By default, the ZeroMQ connection is unauthenticated and unencrypted. Set `HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_CURVE_SERVER_CERT_PATH` to enable [CURVE](https://rfc.zeromq.org/spec/26/), ZeroMQ's built-in encryption and authentication mechanism, for the publisher socket. Both variables are ignored if `ZEROMQ_ENDPOINT` is not set.

| Variable | Required | Default | Description |
|---|---|---|---|
| HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_CURVE_SERVER_CERT_PATH | No — leave empty to disable CURVE |  | Path to the publisher's CZMQ secret cert file (the `*_secret` file produced by `zcert_save` / `goczmq-certgen`), holding the publisher's own CURVE public/secret keypair. Setting this is what turns CURVE on. |
| HOWLITE_RESOURCE_EVENT_PUBLISHER_ZEROMQ_CURVE_ALLOWED_CLIENTS_PATH | No, only relevant if `CURVE_SERVER_CERT_PATH` is set |  | Path to a directory of subscriber public cert files (a CZMQ certstore) allowed to connect. Leave empty to accept any client with a valid CURVE keypair — connections are still encrypted, but not restricted to known peers. |

### Telemetry

All variables are optional. OpenTelemetry export is off by default and turns on automatically as soon as any `OTEL_*` environment variable is set to a non-empty value.

| Variable | Required | Default | Valid values | Description |
|---|---|---|---|---|
| HOWLITE_RESOURCE_TRACING_LEVEL | No | Info | `Debug`, `Info` | Tracing level. Any other/unrecognized value falls back to `Info`. |

Support for standard OTEL environment variables.
Read more [here](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/)

---

## 📄 License

[MIT](LICENSE)
