package main

import (
	"github.com/inx51/howlite/resources/api"
	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/storage/filesystem"
)

func main() {

	storage := filesystem.NewStorage("C:\\test")
	repository := repository.NewRepository(storage)

	api.Run(repository)

}
