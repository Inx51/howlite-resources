package handler_test

import (
	"strings"
	"testing"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/storage/fakestorage"
	"github.com/inx51/howlite/resources/testing/utilities"
	"github.com/inx51/howlite/resources/testing/utilities/tester"
	"github.com/stretchr/testify/assert"
)

func TestCreateShouldReturnCreatedStatusOnSuccess(t *testing.T) {
	tester := tester.NewTester()
	tester.Request.Method = "POST"
	tester.Request.Path = "/test"
	tester.Request.Headers.Set("Content-Type", "application/json") //NOSONAR
	tester.Request.Body.SetString(
		`
		{
			"hello":"world"
		}
		`,
	)

	handler.CreateResource(utilities.CreateHandlerParameters(tester, fakestorage.NewStorage()))
	tester.PopulateResponse()

	assert.Equal(t, 201, tester.Response.StatusCode)
}

func TestCreateShouldReturnLocationHeaderOnSuccess(t *testing.T) {
	tester := tester.NewTester()
	tester.Request.Method = "POST"
	tester.Request.Path = "/test"
	tester.Request.Headers.Set("Content-Type", "application/json")
	tester.Request.Body.SetString(
		`
		{
			"hello":"world"
		}
		`,
	)

	handler.CreateResource(utilities.CreateHandlerParameters(tester, fakestorage.NewStorage()))
	tester.PopulateResponse()

	assert.NotNil(t, tester.Response.Headers.Get("Location"))
	assert.True(t, strings.HasSuffix(tester.Response.Headers.Get("Location"), "/test"))
}

func TestCreateShouldReturnConflictStatusIfResourceExists(t *testing.T) {
	storage := fakestorage.NewStorage()

	storage.AddTestResource("/test", nil, []byte{0x001, 0x002, 0x003})

	tester := tester.NewTester()
	tester.Request.Method = "POST"
	tester.Request.Path = "/test"
	tester.Request.Headers.Set("Content-Type", "application/json")
	tester.Request.Body.SetString(
		`
		{
			"hello":"world"
		}
		`,
	)

	handler.CreateResource(utilities.CreateHandlerParameters(tester, storage))
	tester.PopulateResponse()

	assert.Equal(t, 409, tester.Response.StatusCode)
}
