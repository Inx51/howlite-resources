package main

import (
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage/filesystem"
)

func main() {

	configuration := config.Get()

	storage := filesystem.NewStorage(configuration.Path)
	repository := repository.NewRepository(storage)

	api.Run(repository)

}
