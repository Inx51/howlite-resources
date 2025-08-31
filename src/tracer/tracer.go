package tracer

import (
	"context"
	"strings"

	"github.com/inx51/howlite-resources/configuration"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var level int = -1

func SetupTracer(configuration *configuration.Tracing) {
	tracer = otel.Tracer("howlite-resources")
	setLevel(configuration.LEVEL)
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}

func StartInfoSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if level <= 1 {
		return StartSpan(ctx, name, opts...)
	}
	return ctx, nil
}

func StartDebugSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if level == 0 {
		return StartSpan(ctx, name, opts...)
	}
	return ctx, nil
}

func SafeEndSpan(span trace.Span) {
	if span != nil {
		(span).End()
	}
}

func setLevel(traceLevel string) int {
	if traceLevel == "" {
		traceLevel = "info"
	} else {
		traceLevel = strings.ToLower(traceLevel)
	}

	switch traceLevel {
	case "debug":
		level = 0
	case "info":
		level = 1
	default:
		level = 1
	}

	return level
}

func SetDebugAttributes(ctx context.Context, span trace.Span, kv ...attribute.KeyValue) {
	if level == 0 {
		span.SetAttributes(kv...)
	}
}

func SetInfoAttributes(ctx context.Context, span trace.Span, kv ...attribute.KeyValue) {
	if level <= 1 {
		span.SetAttributes(kv...)
	}
}
