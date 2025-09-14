package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.Printf(
			"[%s] %s %s %d %d %v %s %s",
			start.Format("2006-01-02 15:04:05"),
			r.RemoteAddr,
			r.Method,
			wrapped.statusCode,
			wrapped.size,
			duration,
			r.URL.Path,
			r.UserAgent(),
		)
	}
}

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := `{
					"error": "Internal Server Error",
					"message": "Something went wrong",
					"timestamp": "` + time.Now().Format(time.RFC3339) + `"
				}`
				w.Write([]byte(response))
			}
		}()

		next.ServeHTTP(w, r)
	}
}