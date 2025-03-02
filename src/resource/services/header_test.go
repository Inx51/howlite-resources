package services_test

import (
	"log/slog"
	"testing"

	"github.com/inx51/howlite/resources/resource/services"
	"github.com/inx51/howlite/resources/testing/utilities"
	"github.com/stretchr/testify/assert"
)

func TestShouldFilterOutInvalidResponseHeaders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&utilities.TestingLogWriter{}, nil))
	expected := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Custom":     {"custom-value"},
	}
	destination := map[string][]string{
		"host":            {"localhost"},
		"accept-encoding": {"gzip"},
		"user-agent":      {"test-agent"},
		"Connection":      {"keep-alive"},
		"Accepts":         {"text/html"},
		"Authorization":   {"Bearer token"},
	}

	for k, v := range expected {
		destination[k] = v
	}

	filtered := services.FilterForValidResponseHeaders(&destination, logger)

	assert.Equal(t, expected, *filtered)
}

func TestShouldAllowHeadersThatsNotInvalid(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&utilities.TestingLogWriter{}, nil))
	testHeaders := map[string][]string{
		"X-MyHeader":     {"abcd"},
		"Content-Length": {"123"},
	}

	expected := map[string][]string{
		"X-MyHeader":     {"abcd"},
		"Content-Length": {"123"},
	}

	filtered := services.FilterForValidResponseHeaders(&testHeaders, logger)

	assert.Equal(t, expected, *filtered)
}

func TestShouldPassIfAllProvidedHeadersAreInvalid(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&utilities.TestingLogWriter{}, nil))
	testHeaders := map[string][]string{
		"host":            {"127.0.0.1"},
		"accept-encoding": {"deflate"},
		"connection":      {"close"},
		"Accepts":         {"image/png"},
		"user-agent":      {"agent"},
		"Authorization":   {"Bearer token"},
	}

	expected := map[string][]string{}

	filtered := services.FilterForValidResponseHeaders(&testHeaders, logger)

	assert.Equal(t, expected, *filtered)
}
