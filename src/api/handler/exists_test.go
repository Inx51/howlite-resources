package handler

import (
	"testing"

	"github.com/inx51/howlite/resources/testing/fakes"
	"github.com/inx51/howlite/resources/testing/utilities"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	"github.com/stretchr/testify/assert"
)

func TestExistsShouldReturnNoContentStatusIfResourceExists(t *testing.T) {
	storage := fakes.NewStorage()
	storage.AddTestResource("/test", nil, nil)

	tester := tester.NewTester()
	tester.Request.Method = "HEAD"
	tester.Request.Path = "/test"

	ResourceExists(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 204, tester.Response.StatusCode)
}

func TestExistsShouldReturnCustomHeadersIfOKStatus(t *testing.T) {
	storage := fakes.NewStorage()

	headers := make(map[string][]string)
	headers["X-Custom-Header"] = []string{"custom-value"}

	storage.AddTestResource("/test", headers, nil)

	tester := tester.NewTester()
	tester.Request.Method = "HEAD"
	tester.Request.Path = "/test"

	ResourceExists(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, "custom-value", tester.Response.Headers.Get("X-Custom-Header"))
}

func TestExistsShouldReturnNotFoundStatusIfResourceDoesNotExist(t *testing.T) {
	storage := fakes.NewStorage()

	tester := tester.NewTester()
	tester.Request.Method = "HEAD"
	tester.Request.Path = "/test"

	ResourceExists(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 404, tester.Response.StatusCode)
}
