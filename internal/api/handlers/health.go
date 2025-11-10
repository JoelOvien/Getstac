package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joelovien/go-xlsx-api/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Handle processes the health check request
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
