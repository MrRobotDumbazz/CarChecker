package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Prediction struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	ImageID               uuid.UUID       `json:"image_id" db:"image_id"`
	CleanlinessStatus     string          `json:"cleanliness_status" db:"cleanliness_status"`
	CleanlinessConfidence float64         `json:"cleanliness_confidence" db:"cleanliness_confidence"`
	IntegrityStatus       string          `json:"integrity_status" db:"integrity_status"`
	IntegrityConfidence   float64         `json:"integrity_confidence" db:"integrity_confidence"`
	ProcessingTimeMs      int             `json:"processing_time_ms" db:"processing_time_ms"`
	MLModelVersion        string          `json:"ml_model_version" db:"ml_model_version"`
	AdditionalData        json.RawMessage `json:"additional_data" db:"additional_data"`
	Status                string          `json:"status" db:"status"`
	ErrorMessage          string          `json:"error_message" db:"error_message"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	CompletedAt           *time.Time      `json:"completed_at" db:"completed_at"`
}

type PredictionRequest struct {
	ImageID uuid.UUID `json:"image_id" binding:"required"`
}

type PredictionResponse struct {
	ID          uuid.UUID `json:"id"`
	ImageID     uuid.UUID `json:"image_id"`
	Cleanliness struct {
		Status     string  `json:"status"`
		Confidence float64 `json:"confidence"`
	} `json:"cleanliness"`
	Integrity struct {
		Status     string  `json:"status"`
		Confidence float64 `json:"confidence"`
	} `json:"integrity"`
	ProcessingTimeMs int       `json:"processing_time_ms"`
	ModelVersion     string    `json:"model_version"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Message          string    `json:"message,omitempty"`
}

type MLPredictionRequest struct {
	ImagePath    string `json:"image_path"`
	ModelVersion string `json:"model_version,omitempty"`
}

type MLPredictionResponse struct {
	Cleanliness struct {
		Status     string  `json:"status"`
		Confidence float64 `json:"confidence"`
	} `json:"cleanliness"`
	Integrity struct {
		Status     string  `json:"status"`
		Confidence float64 `json:"confidence"`
	} `json:"integrity"`
	ProcessingTime int    `json:"processing_time_ms"`
	ModelVersion   string `json:"model_version"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
}