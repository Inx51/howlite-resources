package config

type Configuration struct {
	Path string `env:"HOWLITE_RESOURCE_PATH"`
	Host string `env:"HOWLITE_RESOURCE_HOST"`
	Port int    `env:"HOWLITE_RESOURCE_PORT"`
}
