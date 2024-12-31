package api

import (
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/config"
	"github.com/inx51/howlite/resources/api/handler"
)

func Run() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handler.GetResource(&resp, req)
		case "POST":
			handler.CreateResource(&resp, req)
		case "PUT":
			handler.ReplaceResource(&resp, req)
		case "DELETE":
			handler.RemoveResource(&resp, req)
		case "HEAD":
			handler.ResourceExists(&resp, req)
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(config.Instance.HttpServer.Port), nil)
}
