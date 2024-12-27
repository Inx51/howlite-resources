package api

import (
	"net/http"

	"github.com/inx51/howlite/resources/api/handlers"
)

func Run() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handlers.GetResource(&resp, req)
		case "POST":
			handlers.CreateResource(&resp, req)
		case "PUT":
			handlers.ReplaceResource(&resp, req)
		case "DELETE":
			handlers.RemoveResource(&resp, req)
		case "HEAD":
			handlers.ResourceExists(&resp, req)
		}
	})

	http.ListenAndServe(":8080", nil)
}
