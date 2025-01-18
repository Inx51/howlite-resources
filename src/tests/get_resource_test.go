package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/inx51/howlite/resources/tests/fakes"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetResourceShouldReturnNotFoundIfResourceMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://test.local/myresource", nil)
	respRecorder := httptest.NewRecorder()
	storage := fakes.NewStorage()
	handler.GetResource(respRecorder, req, &storage)

	resp := respRecorder.Result()

	assert.Equal(t, resp.StatusCode, 404)
}
