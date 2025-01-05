package config

import (
	"encoding/json"
	"os"

	"github.com/inx51/howlite/resources/path"
)

type Config struct {
	Storage struct {
		Path string `json:"path"`
	}
	HttpServer struct {
		Port int `json:"port"`
	} `json:"httpServer"`
}

type jsonFile struct {
	path     string
	optional bool
}

var Instance *Config
var jsonFiles = []jsonFile{}

func AddJsonConfig(path string, optional bool) {
	jsonFiles = append(jsonFiles, jsonFile{path: path, optional: optional})
}

func Build() {
	buildJsonConfig()
}

func buildJsonConfig() {
	for _, jsonFile := range jsonFiles {
		jsonFile.path = path.Relative(jsonFile.path)
		file, err := os.Open(jsonFile.path)
		if !jsonFile.optional && err != nil {
			panic(err)
		}
		jsonParser := json.NewDecoder(file)
		jsonParser.Decode(&Instance)
	}
}
