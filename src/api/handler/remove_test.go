package handler

import (
	"testing"

	"github.com/inx51/howlite/resources/testing/fakes"
	"github.com/inx51/howlite/resources/testing/utilities"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	"github.com/stretchr/testify/assert"
)

func TestRemoveShouldReturnNoContentStatusIfResourceRemoved(t *testing.T) {
	storage := fakes.NewStorage()
	storage.AddTestResource("/test", nil, nil)

	tester := tester.NewTester()
	tester.Request.Method = "DELETE"
	tester.Request.Path = "/test"

	RemoveResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 204, tester.Response.StatusCode)
}

func TestRemoveShouldReturnNotFoundStatusIfResourceDoesNotExist(t *testing.T) {
	storage := fakes.NewStorage()

	tester := tester.NewTester()
	tester.Request.Method = "HEAD"
	tester.Request.Path = "/test"

	GetResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 404, tester.Response.StatusCode)
}
