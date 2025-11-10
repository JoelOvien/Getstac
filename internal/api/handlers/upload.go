package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/joelovien/go-xlsx-api/internal/storage"
	"github.com/joelovien/go-xlsx-api/internal/xlsx"
	"github.com/rs/zerolog"
)

// UploadHandler handles XLSX file uploads
type UploadHandler struct {
	storage        *storage.MemoryStorage
	parser         *xlsx.Parser
	maxUploadBytes int64
	logger         *zerolog.Logger
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage *storage.MemoryStorage, parser *xlsx.Parser, maxUploadMB int64, logger *zerolog.Logger) *UploadHandler {
	return &UploadHandler{
		storage:        storage,
		parser:         parser,
		maxUploadBytes: maxUploadMB * 1024 * 1024,
		logger:         logger,
	}
}

// Handle processes the upload request
func (h *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Enforce max upload size
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadBytes)

	// Parse multipart form
	err := r.ParseMultipartForm(h.maxUploadBytes)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse multipart form")
		h.writeError(w, http.StatusBadRequest, "bad_request", "File size exceeds maximum allowed size")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get file from form")
		h.writeError(w, http.StatusBadRequest, "bad_request", "Missing or invalid file field")
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" {
		h.logger.Warn().Str("filename", header.Filename).Str("ext", ext).Msg("Invalid file extension")
		h.writeError(w, http.StatusBadRequest, "invalid_file_type", "Only .xlsx files are accepted")
		return
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "spreadsheet") && !strings.Contains(contentType, "excel") && !strings.Contains(contentType, "octet-stream") {
		h.logger.Warn().Str("content_type", contentType).Msg("Invalid content type")
		h.writeError(w, http.StatusBadRequest, "invalid_content_type", "Invalid content type for XLSX file")
		return
	}

	// Generate upload ID
	uploadID := uuid.New().String()

	h.logger.Info().
		Str("upload_id", uploadID).
		Str("filename", header.Filename).
		Int64("size", header.Size).
		Msg("Processing file upload")

	// Read file into memory for parsing
	// Note: For very large files, consider using a temporary file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read file")
		h.writeError(w, http.StatusInternalServerError, "internal_error", "Failed to read uploaded file")
		return
	}

	// Create a reader from the bytes
	reader := strings.NewReader(string(fileBytes))

	// Parse the XLSX file
	result, err := h.parser.Parse(ctx, reader, uploadID)
	if err != nil {
		h.logger.Error().Err(err).Str("upload_id", uploadID).Msg("Failed to parse XLSX file")

		// Provide more specific error messages
		errMsg := err.Error()
		if strings.Contains(errMsg, "no sheets") {
			h.writeError(w, http.StatusBadRequest, "invalid_file", "XLSX file has no sheets")
		} else if strings.Contains(errMsg, "no data") {
			h.writeError(w, http.StatusBadRequest, "invalid_file", "XLSX file has no data")
		} else if strings.Contains(errMsg, "header") {
			h.writeError(w, http.StatusBadRequest, "invalid_headers", errMsg)
		} else {
			h.writeError(w, http.StatusBadRequest, "parse_error", "Failed to parse XLSX file: "+errMsg)
		}
		return
	}

	// Store the parsed records
	if len(result.Records) > 0 {
		err = h.storage.Store(result.Records)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to store records")
			h.writeError(w, http.StatusInternalServerError, "internal_error", "Failed to store records")
			return
		}
	}

	h.logger.Info().
		Str("upload_id", uploadID).
		Int("rows_accepted", result.RowsAccepted).
		Int("rows_rejected", result.RowsRejected).
		Msg("Upload processed successfully")

	// Return response
	response := models.UploadResponse{
		UploadID:     uploadID,
		RowsAccepted: result.RowsAccepted,
		RowsRejected: result.RowsRejected,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeError writes a JSON error response
func (h *UploadHandler) writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
