package services_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/inx51/howlite/resources/resource/services"
	"github.com/inx51/howlite/resources/testing/utilities/logging"
	"github.com/stretchr/testify/assert"
)

func TestShouldFilterOutInvalidResponseHeaders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
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
	ctx := context.Background()

	filtered := services.FilterForValidResponseHeadersContext(ctx, &destination, logger)

	assert.Equal(t, expected, *filtered)
}

func TestShouldAllowHeadersThatsNotInvalid(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	testHeaders := map[string][]string{
		"X-MyHeader":     {"abcd"},
		"Content-Length": {"123"},
	}

	expected := map[string][]string{
		"X-MyHeader":     {"abcd"},
		"Content-Length": {"123"},
	}
	ctx := context.Background()

	filtered := services.FilterForValidResponseHeadersContext(ctx, &testHeaders, logger)

	assert.Equal(t, expected, *filtered)
}

func TestShouldPassIfAllProvidedHeadersAreInvalid(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&logging.TestingLogWriter{}, nil))
	testHeaders := map[string][]string{
		"host":            {"127.0.0.1"},
		"accept-encoding": {"deflate"},
		"connection":      {"close"},
		"Accepts":         {"image/png"},
		"user-agent":      {"agent"},
		"Authorization":   {"Bearer token"},
	}
	ctx := context.Background()

	expected := map[string][]string{}

	filtered := services.FilterForValidResponseHeadersContext(ctx, &testHeaders, logger)

	assert.Equal(t, expected, *filtered)
}
