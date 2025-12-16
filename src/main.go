package main

//TODO: * ADD METRICS

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/telemetry"
)

func main() {
	ctx := context.Background()
	application := NewApplication()

	//Configurations
	configurations := application.ConfigureConfigurations(ctx)

	//Telemetry
	telemetry.SetupLogging(ctx)
	telemetry.SetupMetric(ctx)
	telemetry.SetupTracing(ctx, &configurations.TRACING)
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	application.ConfigureContainer(ctx)
	application.Run(ctx)
	logger.Info(ctx, "Succesfully started application")

	<-ctx.Done()
	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	application.Shutdown(shutdownContext)
	logger.Info(shutdownContext, "Application shutdown gracefully")
	telemetry.ShutdownMetrics(shutdownContext)
	telemetry.ShutdownTracing(shutdownContext)
	telemetry.ShutdownLogging(shutdownContext)
}
