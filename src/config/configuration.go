package config

import "github.com/caarlos0/env/v11"

type Configuration struct {
	Path string `env:"HOWLITE_RESOURCE_PATH"`
}

var instance *Configuration

func (config *Configuration) new() *Configuration {
	instance := Configuration{}
	env.Parse(instance)
	return &instance
}

func Get() *Configuration {
	if instance == nil {
		instance = instance.new()
	}
	return instance
}
