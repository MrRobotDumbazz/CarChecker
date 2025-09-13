package middleware

import (
	"encoding/json"
	"net/http"
)

func JSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}

func WriteErrorWithDetails(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"error":   http.StatusText(status),
		"message": message,
		"details": details,
	}
	json.NewEncoder(w).Encode(response)
}