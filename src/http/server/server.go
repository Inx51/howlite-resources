package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/inx51/howlite-resources/http/handlers"
	"github.com/inx51/howlite-resources/logger"
	"github.com/inx51/howlite-resources/tracer"
	"go.opentelemetry.io/otel/attribute"
)

type Server struct {
	mux        *http.ServeMux
	httpServer *http.Server
}

type TimeoutConfigurations struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewServer(
	host string,
	port int,
	handlers *[]handlers.Handler,
	writeTimout time.Duration,
	readTimeout time.Duration,
	idleTimeout time.Duration) *Server {

	mux := http.NewServeMux()
	for _, handler := range *handlers {
		registerHandler(mux, handler)
	}

	addr := host + ":" + strconv.Itoa(port)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimout,
		IdleTimeout:  idleTimeout,
	}

	server := &Server{
		mux:        mux,
		httpServer: httpServer,
	}

	return server
}

func registerHandler(mux *http.ServeMux, handler handlers.Handler) {
	path := handler.Method() + " " + handler.Path()
	mux.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		ctx, span := tracer.StartInfoSpan(ctx, handler.Method()+" "+request.URL.Path)
		defer tracer.SafeEndSpan(span)
		logger.Debug(ctx, path, "method", request.Method, "path", request.URL.Path)
		statusCode, _ := handler.Handle(ctx, request, response)

		tracer.SetInfoAttributes(ctx,
			span,
			attribute.String("method", request.Method),
			attribute.String("path", request.URL.Path),
			attribute.Int("status", statusCode),
		)
	})
}

func (server *Server) Start(ctx context.Context) {
	logger.Info(ctx, "Starting HTTP server", "address", server.httpServer.Addr)
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			logger.Error(ctx, "Failed to start HTTP server", "error", err)
		}
	}()
}

func (server *Server) Shutdown(ctx context.Context) {
	if err := server.httpServer.Shutdown(ctx); err != nil {
		logger.Error(ctx, "Failed to shutdown HTTP server gracefully", "error", err)
	}
	logger.Info(ctx, "HTTP server shutdown gracefully")
}
