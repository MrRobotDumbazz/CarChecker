package services

import (
	"car-status-backend/internal/database"
	"car-status-backend/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ImageService struct {
	db         *database.DB
	uploadPath string
	maxSize    int64
	allowedTypes []string
}

func NewImageService(db *database.DB, uploadPath string, maxSize int64, allowedTypes []string) *ImageService {
	os.MkdirAll(uploadPath, 0755)

	return &ImageService{
		db:           db,
		uploadPath:   uploadPath,
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
	}
}

func (s *ImageService) UploadImage(file multipart.File, header *multipart.FileHeader) (*models.CarImage, error) {
	if header.Size > s.maxSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, s.maxSize)
	}

	mimeType := header.Header.Get("Content-Type")
	if !s.isAllowedType(mimeType) {
		return nil, fmt.Errorf("file type %s is not allowed", mimeType)
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = s.getExtensionFromMimeType(mimeType)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(s.uploadPath, filename)

	// Get absolute path for ML service
	absolutePath, _ := filepath.Abs(filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	carImage := &models.CarImage{
		ID:           uuid.New(),
		Filename:     filename,
		OriginalName: header.Filename,
		FilePath:     absolutePath,
		FileSize:     header.Size,
		MimeType:     mimeType,
		UploadedAt:   time.Now(),
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO car_images (id, filename, original_name, file_path, file_size, mime_type, uploaded_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = s.db.Exec(query,
		carImage.ID,
		carImage.Filename,
		carImage.OriginalName,
		carImage.FilePath,
		carImage.FileSize,
		carImage.MimeType,
		carImage.UploadedAt,
		carImage.CreatedAt,
	)

	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	return carImage, nil
}

func (s *ImageService) GetImageByID(id uuid.UUID) (*models.CarImage, error) {
	var image models.CarImage
	query := `
		SELECT id, filename, original_name, file_path, file_size, mime_type, uploaded_at, created_at
		FROM car_images
		WHERE id = $1
	`

	err := s.db.Get(&image, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return &image, nil
}

func (s *ImageService) DeleteImage(id uuid.UUID) error {
	image, err := s.GetImageByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM car_images WHERE id = $1`
	_, err = s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete image from database: %w", err)
	}

	if err := os.Remove(image.FilePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

func (s *ImageService) isAllowedType(mimeType string) bool {
	for _, allowed := range s.allowedTypes {
		if strings.EqualFold(mimeType, allowed) {
			return true
		}
	}
	return false
}

func (s *ImageService) getExtensionFromMimeType(mimeType string) string {
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		return ".jpg"
	case "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}