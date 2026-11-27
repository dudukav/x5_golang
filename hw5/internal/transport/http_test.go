package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hw5/internal/model"
	"hw5/internal/service"
	"hw5/internal/storage/memory"
)

func TestRequiredHTTPScenarios(t *testing.T) {
	repository := memory.New()
	router := New(service.New(&repository)).Routes(slog.New(slog.NewTextHandler(io.Discard, nil)))

	created := performRequest(t, router, http.MethodPost, "/api/v1/orders", `{
		"items":[
			{"name":"Keyboard","quantity":2,"price":500000},
			{"name":"Mouse","quantity":1,"price":100000}
		]
	}`)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	var order model.Order
	require.NoError(t, json.Unmarshal(created.Body.Bytes(), &order))
	assert.Equal(t, int64(1100000), order.TotalPrice)
	assert.Equal(t, model.StatusNew, order.Status)
	require.Len(t, order.History, 1)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{name: "list", method: http.MethodGet, path: "/api/v1/orders", wantStatus: http.StatusOK},
		{name: "get", method: http.MethodGet, path: "/api/v1/orders/1", wantStatus: http.StatusOK},
		{name: "missing", method: http.MethodGet, path: "/api/v1/orders/99", wantStatus: http.StatusNotFound, wantCode: "ORDER_NOT_FOUND"},
		{name: "bad id", method: http.MethodGet, path: "/api/v1/orders/nope", wantStatus: http.StatusBadRequest, wantCode: "INVALID_REQUEST"},
		{name: "empty items", method: http.MethodPost, path: "/api/v1/orders", body: `{"items":[]}`, wantStatus: http.StatusBadRequest, wantCode: "INVALID_REQUEST"},
		{name: "malformed json", method: http.MethodPost, path: "/api/v1/orders", body: `{"items":`, wantStatus: http.StatusBadRequest, wantCode: "INVALID_REQUEST"},
		{name: "unknown field", method: http.MethodPost, path: "/api/v1/orders", body: `{"items":[],"secret":true}`, wantStatus: http.StatusBadRequest, wantCode: "INVALID_REQUEST"},
		{name: "status to processing", method: http.MethodPatch, path: "/api/v1/orders/1/status", body: `{"status":"processing"}`, wantStatus: http.StatusOK},
		{name: "invalid transition", method: http.MethodPatch, path: "/api/v1/orders/1/status", body: `{"status":"new"}`, wantStatus: http.StatusConflict, wantCode: "INVALID_STATUS_TRANSITION"},
		{name: "history", method: http.MethodGet, path: "/api/v1/orders/1/history", wantStatus: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performRequest(t, router, test.method, test.path, test.body)
			require.Equal(t, test.wantStatus, response.Code, response.Body.String())
			assert.Contains(t, response.Header().Get("Content-Type"), "application/json")
			if test.wantCode != "" {
				var apiError Error
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &apiError))
				assert.Equal(t, test.wantCode, apiError.Code)
			}
		})
	}

	historyResponse := performRequest(t, router, http.MethodGet, "/api/v1/orders/1/history", "")
	var history []model.StatusChange
	require.NoError(t, json.Unmarshal(historyResponse.Body.Bytes(), &history))
	require.Len(t, history, 2)
	assert.Equal(t, model.StatusNew, history[0].Status)
	assert.Equal(t, model.StatusProcessing, history[1].Status)
	assert.False(t, history[1].ChangedAt.Before(history[0].ChangedAt))
}

func TestHealthAndSwagger(t *testing.T) {
	repository := memory.New()
	router := New(service.New(&repository)).Routes(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := []struct {
		path        string
		wantStatus  int
		wantContent string
	}{
		{path: "/health", wantStatus: http.StatusOK, wantContent: `"status":"ok"`},
		{path: "/swagger/openapi.yaml", wantStatus: http.StatusOK, wantContent: "openapi: 3.0.3"},
		{path: "/swagger/index.html", wantStatus: http.StatusOK, wantContent: "Swagger UI"},
	}

	for _, test := range tests {
		response := performRequest(t, router, http.MethodGet, test.path, "")
		require.Equal(t, test.wantStatus, response.Code, response.Body.String())
		assert.Contains(t, response.Body.String(), test.wantContent)
	}
}

func TestRecoveryReturnsSafeErrorAndLoggingHasRequiredFields(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := middlewareChain(logger, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("sensitive panic details")
	}))

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	assert.NotContains(t, response.Body.String(), "sensitive panic details")
	assert.Contains(t, response.Body.String(), "INTERNAL_ERROR")
	for _, field := range []string{"method", "path", "status", "duration", "request_id"} {
		assert.Contains(t, logs.String(), field)
	}
}

func middlewareChain(logger *slog.Logger, next http.Handler) http.Handler {
	next = recoveryMiddleware(logger)(next)
	next = loggingMiddleware(logger)(next)
	return next
}

func performRequest(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
