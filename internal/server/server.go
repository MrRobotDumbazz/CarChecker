package server

import (
	"car-status-backend/internal/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: server,
		router:     mux,
	}
}

func (s *Server) RegisterRoutes(handlers *Handlers) {
	// Health endpoints
	s.router.HandleFunc("/api/v1/health", s.withMiddleware(handlers.Health.Health))
	s.router.HandleFunc("/api/v1/health/ready", s.withMiddleware(handlers.Health.ReadinessCheck))
	s.router.HandleFunc("/api/v1/health/live", s.withMiddleware(handlers.Health.LivenessCheck))

	// Image endpoints
	s.router.HandleFunc("/api/v1/images/upload", s.withMiddleware(handlers.Upload.UploadImage))
	s.router.HandleFunc("/api/v1/images/", s.withMiddleware(handlers.Upload.GetImage))

	// Prediction endpoints
	s.router.HandleFunc("/api/v1/predict/", s.withMiddleware(handlers.Prediction.PredictImage))
	s.router.HandleFunc("/api/v1/predictions/", s.withMiddleware(handlers.Prediction.GetPrediction))
	s.router.HandleFunc("/api/v1/predictions/stats", s.withMiddleware(handlers.Prediction.GetPredictionStats))

	// API Documentation endpoints
	s.router.HandleFunc("/api/docs", s.withMiddleware(handlers.Swagger.ApiDocsIndex))
	s.router.HandleFunc("/api/docs/", s.withMiddleware(handlers.Swagger.ApiDocsIndex))
	s.router.HandleFunc("/api/docs/swagger", s.withMiddleware(handlers.Swagger.ServeSwaggerUI))
	s.router.HandleFunc("/api/docs/redoc", s.withMiddleware(handlers.Swagger.ServeSwaggerRedoc))
	s.router.HandleFunc("/api/docs/openapi.yaml", s.withMiddleware(handlers.Swagger.ServeOpenAPISpec))

	// Root endpoint
	s.router.HandleFunc("/", s.withMiddleware(s.handleRoot))
}

func (s *Server) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return middleware.CORSMiddleware(
		middleware.LoggingMiddleware(
			middleware.RecoveryMiddleware(
				middleware.JSONMiddleware(next),
			),
		),
	)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse := `{
		"service": "car-status-backend",
		"version": "1.0.0",
		"status": "running",
		"description": "API для системы определения состояния автомобиля по фотографии",
		"documentation": {
			"swagger_ui": "/api/docs/swagger",
			"redoc": "/api/docs/redoc",
			"openapi_spec": "/api/docs/openapi.yaml",
			"docs_index": "/api/docs"
		},
		"endpoints": {
			"health": "/api/v1/health",
			"upload": "/api/v1/images/upload",
			"get_image": "/api/v1/images/{id}",
			"predict": "/api/v1/predict/{image_id}",
			"get_prediction": "/api/v1/predictions/{id}",
			"prediction_stats": "/api/v1/predictions/stats"
		}
	}`

	w.Write([]byte(jsonResponse))
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}