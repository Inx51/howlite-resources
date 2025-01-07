package api

import (
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/storage"
)

func Run() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		storage := storage.Create()

		switch req.Method {
		case "GET":
			handler.GetResource(resp, req, &storage)
		case "POST":
			handler.CreateResource(resp, req, &storage)
		case "PUT":
			handler.ReplaceResource(resp, req, &storage)
		case "DELETE":
			handler.RemoveResource(resp, req, &storage)
		case "HEAD":
			handler.ResourceExists(resp, req, &storage)
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(config.Instance.HttpServer.Port), nil)
}
