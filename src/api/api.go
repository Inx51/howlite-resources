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
		store := storage.Create()

		switch req.Method {
		case "GET":
			handler.GetResource(resp, req, store)
		case "POST":
			handler.CreateResource(resp, req, store)
		case "PUT":
			handler.ReplaceResource(resp, req, store)
		case "DELETE":
			handler.RemoveResource(resp, req, store)
		case "HEAD":
			handler.ResourceExists(resp, req, store)
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(config.Instance.HttpServer.Port), nil)
}
