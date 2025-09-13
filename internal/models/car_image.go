package models

import (
	"time"

	"github.com/google/uuid"
)

type CarImage struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_name"`
	FilePath     string    `json:"file_path" db:"file_path"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CarImageRequest struct {
	File []byte `json:"file" binding:"required"`
}

type CarImageResponse struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	Message      string    `json:"message,omitempty"`
}