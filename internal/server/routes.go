package server

import (
	"car-status-backend/internal/database"
	"car-status-backend/internal/handlers"
	"car-status-backend/internal/services"
)

type Handlers struct {
	Health     *handlers.HealthHandler
	Upload     *handlers.UploadHandler
	Prediction *handlers.PredictionHandler
	Swagger    *handlers.SwaggerHandler
}

func NewHandlers(
	imageService *services.ImageService,
	predictionService *services.PredictionService,
	mlClient *services.MLClient,
	queueService *services.QueueService,
	db interface{},
) *Handlers {
	return &Handlers{
		Health:     handlers.NewHealthHandler(db.(*database.DB), mlClient),
		Upload:     handlers.NewUploadHandler(imageService),
		Prediction: handlers.NewPredictionHandler(imageService, predictionService, mlClient, queueService),
		Swagger:    handlers.NewSwaggerHandler("./api/openapi.yaml"),
	}
}