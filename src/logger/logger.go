package logger

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
}

var (
	loggers = make([]Logger, 0)
)

func init() {
	if !isTestRun() {
		registerSlog()
		registerOtel()
	}
}

func registerSlog() {
	loggers = append(loggers, slog.New(slog.NewTextHandler(os.Stdout, nil)))
}

func registerOtel() {
	loggers = append(loggers, otelslog.NewLogger("howlite-resources"))
}

func isTestRun() bool {
	return flag.Lookup("test.v") != nil
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	for _, logger := range loggers {
		logger.DebugContext(ctx, msg, args...)
	}
}

func Info(ctx context.Context, msg string, args ...interface{}) {
	for _, logger := range loggers {
		logger.InfoContext(ctx, msg, args...)
	}
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	for _, logger := range loggers {
		logger.ErrorContext(ctx, msg, args...)
	}
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	for _, logger := range loggers {
		logger.WarnContext(ctx, msg, args...)
	}
}
