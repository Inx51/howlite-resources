package config

type Configuration struct {
	PATH             string `env:"HOWLITE_RESOURCE_PATH"`
	HOST             string `env:"HOWLITE_RESOURCE_HOST"`
	PORT             int    `env:"HOWLITE_RESOURCE_PORT"`
	STORAGE_PROVIDER StorageProvider
	OTEL             OtelConfiguration
}

type OtelConfiguration struct {
	OTEL_SERVICE_NAME                   string `env:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_PROTOCOL         string `env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	OTEL_EXPORTER_OTLP_TRACES_PROTOCOL  string `env:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"`
	OTEL_EXPORTER_OTLP_METRICS_PROTOCOL string `env:"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"`
	OTEL_EXPORTER_OTLP_LOGS_PROTOCOL    string `env:"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"`
}

type StorageProvider struct {
	STORAGE_PROVIDER            string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER"`
	STORAGE_PROVIDER_FILESYSTEM FilesystemConfiguration
	STORAGE_PROVIDER_S3         S3Configuration
}

type FilesystemConfiguration struct {
	PATH string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_FILESYSTEM_PATH"`
}

type S3Configuration struct {
	BUCKET                     string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_BUCKET"`
	PREFIX                     string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_PREFIX"`
	ACCESS_KEY                 string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ACCESS_KEY"`
	SECRET_KEY                 string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_SECRET_KEY"`
	ENDPOINT                   string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_ENDPOINT"`
	REGION                     string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_REGION"`
	UPLOAD_STRATEGY            string `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_UPLOAD_STRATEGY" envDefault:"singlepart"`
	MULTIPART_PART_UPLOAD_SIZE int    `env:"HOWLITE_RESOURCE_STORAGE_PROVIDER_S3_MULTIPART_PART_UPLOAD_SIZE" envDefault:"5242880"`
}
