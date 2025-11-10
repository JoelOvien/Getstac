package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/joelovien/go-xlsx-api/internal/storage"
	"github.com/rs/zerolog"
)

// ListHandler handles listing records with pagination
type ListHandler struct {
	storage *storage.MemoryStorage
	logger  *zerolog.Logger
}

// NewListHandler creates a new list handler
func NewListHandler(storage *storage.MemoryStorage, logger *zerolog.Logger) *ListHandler {
	return &ListHandler{
		storage: storage,
		logger:  logger,
	}
}

// Handle processes the list request
func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 10
	offset := 0

	// Parse limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			h.writeError(w, http.StatusBadRequest, "invalid_parameter", "Invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	// Parse offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			h.writeError(w, http.StatusBadRequest, "invalid_parameter", "Invalid offset parameter")
			return
		}
		offset = parsedOffset
	}

	// Enforce maximum limit to prevent abuse
	if limit > 1000 {
		limit = 1000
	}

	// Retrieve records from storage
	records, total, err := h.storage.List(limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list records")
		h.writeError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve records")
		return
	}

	// Build response
	response := models.ListRecordsResponse{
		Records: records,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}

	h.logger.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Int("total", total).
		Int("returned", len(records)).
		Msg("Listed records")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeError writes a JSON error response
func (h *ListHandler) writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
