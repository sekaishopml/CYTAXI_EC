package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ErrorHandler struct {
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (h *ErrorHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, map[string]any{
		"error":  "route not found",
		"path":   r.URL.Path,
		"method": r.Method,
	})
}

func (h *ErrorHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
		"error":  "method not allowed",
		"method": r.Method,
	})
}

func (h *ErrorHandler) BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]any{"error": msg})
}

func (h *ErrorHandler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type HealthChecker struct {
	logger *slog.Logger
}

func NewHealthChecker(logger *slog.Logger) *HealthChecker {
	return &HealthChecker{logger: logger}
}

func (hc *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"service":   "api-gateway",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
