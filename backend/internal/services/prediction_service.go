package services

import (
	"car-status-backend/internal/database"
	"car-status-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PredictionService struct {
	db *database.DB
}

func NewPredictionService(db *database.DB) *PredictionService {
	return &PredictionService{
		db: db,
	}
}

func (s *PredictionService) CreatePrediction(imageID uuid.UUID, mlResult *models.MLPredictionResponse) (*models.Prediction, error) {
	prediction := &models.Prediction{
		ID:                    uuid.New(),
		ImageID:               imageID,
		ProcessingTimeMs:      mlResult.ProcessingTime,
		MLModelVersion:        mlResult.ModelVersion,
		CreatedAt:            time.Now(),
	}

	if mlResult.Success {
		prediction.CleanlinessStatus = mlResult.Cleanliness.Status
		prediction.CleanlinessConfidence = mlResult.Cleanliness.Confidence
		prediction.IntegrityStatus = mlResult.Integrity.Status
		prediction.IntegrityConfidence = mlResult.Integrity.Confidence
		prediction.Status = "completed"
		now := time.Now()
		prediction.CompletedAt = &now
	} else {
		prediction.Status = "failed"
		prediction.ErrorMessage = mlResult.Error
	}

	additionalData := map[string]interface{}{
		"ml_response": mlResult,
	}
	additionalDataJSON, _ := json.Marshal(additionalData)
	prediction.AdditionalData = additionalDataJSON

	query := `
		INSERT INTO predictions (
			id, image_id, cleanliness_status, cleanliness_confidence,
			integrity_status, integrity_confidence, processing_time_ms,
			ml_model_version, additional_data, status, error_message,
			created_at, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := s.db.Exec(query,
		prediction.ID,
		prediction.ImageID,
		prediction.CleanlinessStatus,
		prediction.CleanlinessConfidence,
		prediction.IntegrityStatus,
		prediction.IntegrityConfidence,
		prediction.ProcessingTimeMs,
		prediction.MLModelVersion,
		prediction.AdditionalData,
		prediction.Status,
		prediction.ErrorMessage,
		prediction.CreatedAt,
		prediction.CompletedAt,
	)

	if err != nil {
		fmt.Printf("DEBUG: Error saving prediction: %v\n", err)
		fmt.Printf("DEBUG: CleanlinessStatus: %v\n", prediction.CleanlinessStatus)
		fmt.Printf("DEBUG: IntegrityStatus: %v\n", prediction.IntegrityStatus)
		fmt.Printf("DEBUG: Status: %v\n", prediction.Status)
		return nil, fmt.Errorf("failed to save prediction: %w", err)
	}

	return prediction, nil
}

func (s *PredictionService) GetPredictionByID(id uuid.UUID) (*models.Prediction, error) {
	var prediction models.Prediction
	query := `
		SELECT id, image_id, cleanliness_status, cleanliness_confidence,
		       integrity_status, integrity_confidence, processing_time_ms,
		       ml_model_version, additional_data, status, error_message,
		       created_at, completed_at
		FROM predictions
		WHERE id = $1
	`

	err := s.db.Get(&prediction, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	return &prediction, nil
}

func (s *PredictionService) GetPredictionsByImageID(imageID uuid.UUID) ([]models.Prediction, error) {
	var predictions []models.Prediction
	query := `
		SELECT id, image_id, cleanliness_status, cleanliness_confidence,
		       integrity_status, integrity_confidence, processing_time_ms,
		       ml_model_version, additional_data, status, error_message,
		       created_at, completed_at
		FROM predictions
		WHERE image_id = $1
		ORDER BY created_at DESC
	`

	err := s.db.Select(&predictions, query, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}

	return predictions, nil
}

func (s *PredictionService) UpdatePredictionStatus(id uuid.UUID, status string, errorMessage string) error {
	query := `
		UPDATE predictions
		SET status = $1, error_message = $2, completed_at = $3
		WHERE id = $4
	`

	var completedAt *time.Time
	if status == "completed" || status == "failed" {
		now := time.Now()
		completedAt = &now
	}

	_, err := s.db.Exec(query, status, errorMessage, completedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update prediction status: %w", err)
	}

	return nil
}

func (s *PredictionService) GetPredictionStats() (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing,
			AVG(CASE WHEN processing_time_ms > 0 THEN processing_time_ms END) as avg_processing_time
		FROM predictions
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`

	var stats struct {
		Total              int     `db:"total"`
		Completed          int     `db:"completed"`
		Failed             int     `db:"failed"`
		Pending            int     `db:"pending"`
		Processing         int     `db:"processing"`
		AvgProcessingTime  *float64 `db:"avg_processing_time"`
	}

	err := s.db.Get(&stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction stats: %w", err)
	}

	result := map[string]interface{}{
		"total":      stats.Total,
		"completed":  stats.Completed,
		"failed":     stats.Failed,
		"pending":    stats.Pending,
		"processing": stats.Processing,
	}

	if stats.AvgProcessingTime != nil {
		result["avg_processing_time_ms"] = *stats.AvgProcessingTime
	}

	return result, nil
}