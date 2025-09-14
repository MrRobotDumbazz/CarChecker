package handlers

import (
	"car-status-backend/internal/models"
	"car-status-backend/internal/services"
	"car-status-backend/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PredictionHandler struct {
	imageService      *services.ImageService
	predictionService *services.PredictionService
	mlClient          *services.MLClient
	queueService      *services.QueueService
}

func NewPredictionHandler(
	imageService *services.ImageService,
	predictionService *services.PredictionService,
	mlClient *services.MLClient,
	queueService *services.QueueService,
) *PredictionHandler {
	return &PredictionHandler{
		imageService:      imageService,
		predictionService: predictionService,
		mlClient:          mlClient,
		queueService:      queueService,
	}
}

func (h *PredictionHandler) PredictImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	imageIDStr := utils.ExtractIDFromPath(r.URL.Path, "/api/v1/predict/")
	if imageIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	if err := utils.ValidateUUID(imageIDStr); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	imageID, _ := uuid.Parse(imageIDStr)
	image, err := h.imageService.GetImageByID(imageID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Image not found")
		return
	}

	if h.queueService != nil {
		err = h.queueService.PublishPredictionJob(imageID, image.FilePath)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to queue prediction job")
			return
		}

		utils.WriteSuccessResponse(w, http.StatusAccepted, map[string]interface{}{
			"image_id": imageID,
			"status":   "queued",
			"message":  "Prediction job has been queued for processing",
		}, "Prediction job queued successfully")
		return
	}

	start := time.Now()
	mlResult, err := h.mlClient.PredictCarStatus(image.FilePath)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process image with ML service: "+err.Error())
		return
	}

	prediction, err := h.predictionService.CreatePrediction(imageID, mlResult)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to save prediction result")
		return
	}

	response := h.buildPredictionResponse(prediction)
	response.Message = "Prediction completed successfully"

	processingTime := time.Since(start)
	if processingTime.Milliseconds() > 0 {
		response.ProcessingTimeMs = int(processingTime.Milliseconds())
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response, "Prediction completed successfully")
}

func (h *PredictionHandler) GetPrediction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	predictionIDStr := utils.ExtractIDFromPath(r.URL.Path, "/api/v1/predictions/")
	if predictionIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Prediction ID is required")
		return
	}

	if err := utils.ValidateUUID(predictionIDStr); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid prediction ID format")
		return
	}

	predictionID, _ := uuid.Parse(predictionIDStr)
	prediction, err := h.predictionService.GetPredictionByID(predictionID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Prediction not found")
		return
	}

	response := h.buildPredictionResponse(prediction)
	utils.WriteSuccessResponse(w, http.StatusOK, response, "Prediction retrieved successfully")
}

func (h *PredictionHandler) GetImagePredictions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	imageIDStr := utils.ExtractIDFromPath(r.URL.Path, "/api/v1/images/")
	imageIDStr = utils.ExtractIDFromPath(imageIDStr, "predictions")

	if imageIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	if err := utils.ValidateUUID(imageIDStr); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	imageID, _ := uuid.Parse(imageIDStr)
	predictions, err := h.predictionService.GetPredictionsByImageID(imageID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get predictions")
		return
	}

	var responses []models.PredictionResponse
	for _, prediction := range predictions {
		response := h.buildPredictionResponse(&prediction)
		responses = append(responses, response)
	}

	utils.WriteSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"image_id":    imageID,
		"predictions": responses,
		"count":       len(responses),
	}, "Predictions retrieved successfully")
}

func (h *PredictionHandler) GetPredictionStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := h.predictionService.GetPredictionStats()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get prediction stats")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, stats, "Prediction stats retrieved successfully")
}

func (h *PredictionHandler) buildPredictionResponse(prediction *models.Prediction) models.PredictionResponse {
	response := models.PredictionResponse{
		ID:               prediction.ID,
		ImageID:          prediction.ImageID,
		ProcessingTimeMs: prediction.ProcessingTimeMs,
		ModelVersion:     prediction.MLModelVersion,
		Status:           prediction.Status,
		CreatedAt:        prediction.CreatedAt,
		CompletedAt:      prediction.CompletedAt,
	}

	if prediction.Status == "completed" {
		response.Cleanliness.Status = prediction.CleanlinessStatus
		response.Cleanliness.Confidence = prediction.CleanlinessConfidence
		response.Integrity.Status = prediction.IntegrityStatus
		response.Integrity.Confidence = prediction.IntegrityConfidence
	}

	return response
}