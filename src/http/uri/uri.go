package uri

import "net/http"

func AbsoluteUri(request *http.Request) string {
	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}

	return schema + "://" + request.URL.Host + request.URL.Path
}
