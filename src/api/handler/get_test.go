package handler_test

import (
	"testing"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/storage/fakestorage"
	"github.com/inx51/howlite/resources/testing/utilities"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	"github.com/stretchr/testify/assert"
)

func TestGetShouldReturnOKStatusIfResourceExists(t *testing.T) {
	storage := fakestorage.NewStorage()
	storage.AddTestResource("/test", nil, nil)

	tester := tester.NewTester()
	tester.Request.Method = "GET"
	tester.Request.Path = "/test"

	handler.GetResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 200, tester.Response.StatusCode)
}

func TestGetShouldReturnCustomHeadersIfOKStatus(t *testing.T) {
	storage := fakestorage.NewStorage()

	headers := make(map[string][]string)
	headers["X-Custom-Header"] = []string{"custom-value"}

	storage.AddTestResource("/test", headers, nil)

	tester := tester.NewTester()
	tester.Request.Method = "GET"
	tester.Request.Path = "/test"

	handler.GetResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, "custom-value", tester.Response.Headers.Get("X-Custom-Header"))
}

func TestGetShouldReturnBodyIfOKStatus(t *testing.T) {
	storage := fakestorage.NewStorage()

	body := []byte{0x001, 0x002, 0x003}

	storage.AddTestResource("/test", nil, body)

	tester := tester.NewTester()
	tester.Request.Method = "GET"
	tester.Request.Path = "/test"

	handler.GetResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, body, tester.Response.Body)
}

func TestGetShouldReturnNotFoundStatusIfResourceDoesNotExist(t *testing.T) {
	storage := fakestorage.NewStorage()

	tester := tester.NewTester()
	tester.Request.Method = "HEAD"
	tester.Request.Path = "/test"

	handler.GetResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 404, tester.Response.StatusCode)
}
