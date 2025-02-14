package api

import (
	"net/http"
	"strconv"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/resource/repository"
)

func SetupHandlers(repository *repository.Repository) {
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
}

func Run(host string, port int) {
	http.ListenAndServe(host+":"+strconv.Itoa(port), nil)
}
