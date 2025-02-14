package tester

import (
	"net/http"
	"net/http/httptest"
)

type Tester struct {
	Request          Request
	Response         Response
	responseRecorder *httptest.ResponseRecorder
}

func NewTester() *Tester {
	return &Tester{
		Request: NewRequest(),
	}
}

func (tester *Tester) Build() (http.ResponseWriter, *http.Request) {
	request := httptest.NewRequest(tester.Request.Method, tester.Request.Path, tester.Request.Body.reader)
	recorder := httptest.NewRecorder()
	tester.responseRecorder = recorder
	return recorder, request
}

func (tester *Tester) PopulateResponse() {
	response := tester.responseRecorder.Result()
	(*tester).Response = NewResponse(*response)
}
