package transport

import (
	"embed"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"hw5/internal/model"
	"hw5/internal/service"
)

//go:embed openapi.yaml
var swaggerFiles embed.FS

func (h *Server) Routes(logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(loggingMiddleware(logger))
	router.Use(recoveryMiddleware(logger))

	strictHandler := NewStrictHandlerWithOptions(h, nil, StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			writeAPIError(w, wrapInvalidRequest(err))
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			writeAPIError(w, err)
		},
	})
	HandlerWithOptions(strictHandler, ChiServerOptions{
		BaseRouter: router,
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			writeAPIError(w, wrapInvalidRequest(err))
		},
	})

	router.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusTemporaryRedirect)
	})
	router.Get("/swagger/openapi.yaml", swaggerSpec)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/openapi.yaml"),
		httpSwagger.DocExpansion("list"),
	))

	return router
}

func swaggerSpec(w http.ResponseWriter, _ *http.Request) {
	spec, err := swaggerFiles.ReadFile("openapi.yaml")
	if err != nil {
		writeAPIError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(spec)
}

func wrapInvalidRequest(err error) error {
	return errors.Join(service.ErrInvalidOrder, err)
}

func writeAPIError(w http.ResponseWriter, err error) {
	response := Error{Error: "unexpected internal error", Code: "INTERNAL_ERROR"}
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, service.ErrInvalidOrder):
		status = http.StatusBadRequest
		response = Error{Error: err.Error(), Code: "INVALID_REQUEST"}
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
		response = Error{Error: "order not found", Code: "ORDER_NOT_FOUND"}
	case errors.Is(err, service.ErrInvalidStatusTransition):
		status = http.StatusConflict
		response = Error{Error: "status transition is not allowed", Code: "INVALID_STATUS_TRANSITION"}
	}
	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (w *responseRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseRecorder) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &responseRecorder{ResponseWriter: w}
			next.ServeHTTP(recorder, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.status,
				"duration", time.Since(started),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}

func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						"panic", recovered,
						"request_id", middleware.GetReqID(r.Context()),
					)
					writeJSON(w, http.StatusInternalServerError, Error{
						Error: "unexpected internal error",
						Code:  "INTERNAL_ERROR",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
