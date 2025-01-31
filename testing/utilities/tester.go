package utilities

type Tester struct {
	Request  Request
	Response Response
}

type Request struct {
	Path    string
	Method  string
	Headers Headers
	Body    interface{}
}

type Response struct {
	Status  uint
	Headers Headers
	Body    interface{}
}

type Headers struct {
	headers map[string]string
}

func (headers *Headers) Add(key string, value string) {
	headers.headers[key] = value
}

func NewTester() *Tester {
	return &Tester{}
}
