package handlers

import (
	"car-status-backend/internal/database"
	"car-status-backend/internal/services"
	"car-status-backend/pkg/utils"
	"net/http"
)

type HealthHandler struct {
	db       *database.DB
	mlClient *services.MLClient
}

func NewHealthHandler(db *database.DB, mlClient *services.MLClient) *HealthHandler {
	return &HealthHandler{
		db:       db,
		mlClient: mlClient,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	checks := make(map[string]string)

	if err := h.db.HealthCheck(); err != nil {
		checks["database"] = "failed: " + err.Error()
	} else {
		checks["database"] = "ok"
	}

	if err := h.mlClient.HealthCheck(); err != nil {
		checks["ml_service"] = "failed: " + err.Error()
	} else {
		checks["ml_service"] = "ok"
	}

	utils.WriteHealthResponse(w, "car-status-backend", checks)
}

func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := h.db.HealthCheck(); err != nil {
		utils.WriteErrorResponse(w, http.StatusServiceUnavailable, "Database not ready")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, map[string]string{
		"status": "ready",
	}, "Service is ready")
}

func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, map[string]string{
		"status": "alive",
	}, "Service is alive")
}