package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type ErrorResponse struct {
	Error     string      `json:"error"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteSuccessResponse(w http.ResponseWriter, status int, data interface{}, message string) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
	WriteJSON(w, status, response)
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	response := ErrorResponse{
		Error:     http.StatusText(status),
		Message:   message,
		Timestamp: time.Now(),
	}
	WriteJSON(w, status, response)
}

func WriteErrorResponseWithDetails(w http.ResponseWriter, status int, message string, details interface{}) {
	response := ErrorResponse{
		Error:     http.StatusText(status),
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
	WriteJSON(w, status, response)
}

func WriteValidationError(w http.ResponseWriter, errors map[string]string) {
	WriteErrorResponseWithDetails(w, http.StatusBadRequest, "Validation failed", errors)
}

func WriteHealthResponse(w http.ResponseWriter, service string, checks map[string]string) {
	status := "healthy"
	httpStatus := http.StatusOK

	if checks != nil {
		for _, check := range checks {
			if check != "ok" {
				status = "unhealthy"
				httpStatus = http.StatusServiceUnavailable
				break
			}
		}
	}

	response := HealthResponse{
		Status:    status,
		Service:   service,
		Timestamp: time.Now(),
		Checks:    checks,
	}

	WriteJSON(w, httpStatus, response)
}

func ParseRequestBody(r *http.Request, dest interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}

func SetJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func SetCacheHeaders(w http.ResponseWriter, maxAge int) {
	if maxAge <= 0 {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	}
}

func ExtractIDFromPath(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimPrefix(path, prefix)
}