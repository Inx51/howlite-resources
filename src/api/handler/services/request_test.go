package services

import (
	"crypto/tls"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestUrlShouldGetRequestUrlAsHttp(t *testing.T) {
	expected := "http://example.com/some/path"
	req := httptest.NewRequest("GET", expected, nil)
	actual := GetRequestUrl(req)

	assert.Equal(t, expected, actual)
}

func TestGetRequestUrlShouldGetRequestUrlAsHttps(t *testing.T) {
	expected := "https://example.com/some/path"
	req := httptest.NewRequest("GET", expected, nil)
	req.TLS = &tls.ConnectionState{}
	actual := GetRequestUrl(req)

	assert.Equal(t, expected, actual)
}
