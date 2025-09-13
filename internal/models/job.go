package models

import (
	"time"

	"github.com/google/uuid"
)

type PredictionJob struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ImageID      uuid.UUID  `json:"image_id" db:"image_id"`
	Status       string     `json:"status" db:"status"`
	RetryCount   int        `json:"retry_count" db:"retry_count"`
	MaxRetries   int        `json:"max_retries" db:"max_retries"`
	ScheduledAt  time.Time  `json:"scheduled_at" db:"scheduled_at"`
	StartedAt    *time.Time `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage string     `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type PredictionJobRequest struct {
	ImageID   uuid.UUID `json:"image_id"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}

type PredictionJobResponse struct {
	ID          uuid.UUID  `json:"id"`
	ImageID     uuid.UUID  `json:"image_id"`
	Status      string     `json:"status"`
	RetryCount  int        `json:"retry_count"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Message     string     `json:"message,omitempty"`
}