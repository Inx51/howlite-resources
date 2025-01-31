package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inx51/howlite/resources/resource/repository"
	"github.com/inx51/howlite/resources/testing/fakes"
	"github.com/inx51/howlite/resources/testing/utilities"
)

func TestShouldReturnCreatedStatusOnSuccess(t *testing.T) {

	resourceBody :=
		`
		{
			"hello":"world"
		}
		`
	body := strings.NewReader(resourceBody)

	request := httptest.NewRequest("GET", "http://localhost:8080/test", body)
	recorder := httptest.NewRecorder()

	storage := fakes.NewStorage()
	repository := repository.NewRepository(storage)

	CreateResource(recorder, request, repository)

	result := recorder.Result()
	if result.StatusCode != 201 {
		t.Errorf("Returned StatusCode of %s, expected 201", result.Status)
	}
}

func TestShouldReturnLocationHeaderOnSuccess(t *testing.T) {

	tester := utilities.NewTester()
	tester.Request.Method = ""
	tester.Request.Path = ""
	tester.Request.Headers
	jsonBody := tester.NewJsonBody()
	jsonBody[""] = ""
	tester.Request.Body.SetJson(jsonBody)

	CreateResource(tester.Build(), repository)

	Is.Equals(tester.Response.Headers["Location"], resourceUrl)

	// resourceBody :=
	// 	`
	// 	{
	// 		"hello":"world"
	// 	}
	// 	`
	// body := strings.NewReader(resourceBody)

	// resourceUrl := "http://localhost:8080/test"
	// request := httptest.NewRequest("GET", resourceUrl, body)
	// recorder := httptest.NewRecorder()

	// storage := fakes.NewStorage()
	// repository := repository.NewRepository(storage)

	// CreateResource(recorder, request, repository)

	// location := recorder.Result().Header["Location"]
	// if location[0] != resourceUrl {
	// 	t.Errorf("Failed to return a Location header")
	// }
}
