package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/inx51/howlite/resources/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetResourceShouldReturnNotFoundIfResourceMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://test.local/myresource", nil)
	respRecorder := httptest.NewRecorder()
	handler.GetResource(respRecorder, req)

	resp := respRecorder.Result()

	assert.Equal(t, resp.StatusCode, 404)
}
