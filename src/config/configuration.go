package config

type Configuration struct {
	PATH string `env:"HOWLITE_RESOURCE_PATH"`
	HOST string `env:"HOWLITE_RESOURCE_HOST"`
	PORT int    `env:"HOWLITE_RESOURCE_PORT"`

	OTEL OtelConfiguration
}

type OtelConfiguration struct {
	OTEL_SERVICE_NAME                   string `env:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_PROTOCOL         string `env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	OTEL_EXPORTER_OTLP_TRACES_PROTOCOL  string `env:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"`
	OTEL_EXPORTER_OTLP_METRICS_PROTOCOL string `env:"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"`
	OTEL_EXPORTER_OTLP_LOGS_PROTOCOL    string `env:"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"`
	DEV_MODE                            bool   `env:"DEV_MODE"`
}
