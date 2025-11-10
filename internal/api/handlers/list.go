package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/joelovien/go-xlsx-api/internal/storage"
	"github.com/rs/zerolog"
)

type ListHandler struct {
	storage *storage.MemoryStorage
	logger  *zerolog.Logger
}

func NewListHandler(storage *storage.MemoryStorage, logger *zerolog.Logger) *ListHandler {
	return &ListHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			h.writeError(w, http.StatusBadRequest, "invalid_parameter", "Invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			h.writeError(w, http.StatusBadRequest, "invalid_parameter", "Invalid offset parameter")
			return
		}
		offset = parsedOffset
	}

	if limit > 1000 {
		limit = 1000
	}

	records, total, err := h.storage.List(limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list records")
		h.writeError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve records")
		return
	}

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

func (h *ListHandler) writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
