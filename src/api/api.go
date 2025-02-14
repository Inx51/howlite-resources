package api

import (
	"net/http"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
)

func Run(repository *repository.Repository) {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handler.GetResource(resp, req, repository)
		case "POST":
			handler.CreateResource(resp, req, repository)
		case "PUT":
			handler.ReplaceResource(resp, req, repository)
		case "DELETE":
			handler.RemoveResource(resp, req, repository)
		case "HEAD":
			handler.ResourceExists(resp, req, repository)
		default:
			resp.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}
