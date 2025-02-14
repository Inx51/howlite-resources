package tester

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

func (headers *Headers) Set(key string, value string) {
	headers.headers[key] = value
}

func (headers *Headers) Get(key string) string {
	return headers.headers[key]
}
