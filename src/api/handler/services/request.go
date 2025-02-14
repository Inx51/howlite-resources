package services

import "net/http"

func GetRequestUrl(request *http.Request) string {
	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}

	return schema + "://" + request.URL.Host + request.URL.Path
}
