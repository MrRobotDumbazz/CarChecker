package handlers

import (
	"car-status-backend/internal/models"
	"car-status-backend/internal/services"
	"car-status-backend/pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

type UploadHandler struct {
	imageService *services.ImageService
}

func NewUploadHandler(imageService *services.ImageService) *UploadHandler {
	return &UploadHandler{
		imageService: imageService,
	}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	if header == nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid file header")
		return
	}

	image, err := h.imageService.UploadImage(file, header)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := models.CarImageResponse{
		ID:           image.ID,
		Filename:     image.Filename,
		OriginalName: image.OriginalName,
		FileSize:     image.FileSize,
		MimeType:     image.MimeType,
		UploadedAt:   image.UploadedAt,
		Message:      "Image uploaded successfully",
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, response, "Image uploaded successfully")
}

func (h *UploadHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	imageIDStr := utils.ExtractIDFromPath(r.URL.Path, "/api/v1/images/")
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

	response := models.CarImageResponse{
		ID:           image.ID,
		Filename:     image.Filename,
		OriginalName: image.OriginalName,
		FileSize:     image.FileSize,
		MimeType:     image.MimeType,
		UploadedAt:   image.UploadedAt,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, response, "Image retrieved successfully")
}

func (h *UploadHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	imageIDStr := utils.ExtractIDFromPath(r.URL.Path, "/api/v1/images/")
	if imageIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	if err := utils.ValidateUUID(imageIDStr); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	imageID, _ := uuid.Parse(imageIDStr)
	err := h.imageService.DeleteImage(imageID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, map[string]string{
		"id": imageID.String(),
	}, "Image deleted successfully")
}