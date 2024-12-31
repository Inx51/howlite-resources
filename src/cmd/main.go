package main

import (
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/api/config"
)

func main() {
	config.AddJsonConfig("./config.json", false)
	config.AddJsonConfig("./config.secret.json", true)
	config.Build()
	api.Run()
}
