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

type UploadHandler struct {
	storage        *storage.MemoryStorage
	parser         *xlsx.Parser
	maxUploadBytes int64
	logger         *zerolog.Logger
}

func NewUploadHandler(storage *storage.MemoryStorage, parser *xlsx.Parser, maxUploadMB int64, logger *zerolog.Logger) *UploadHandler {
	return &UploadHandler{
		storage:        storage,
		parser:         parser,
		maxUploadBytes: maxUploadMB * 1024 * 1024,
		logger:         logger,
	}
}

func (h *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadBytes)

	err := r.ParseMultipartForm(h.maxUploadBytes)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse multipart form")
		h.writeError(w, http.StatusBadRequest, "bad_request", "File size exceeds maximum allowed size")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get file from form")
		h.writeError(w, http.StatusBadRequest, "bad_request", "Missing or invalid file field")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" {
		h.logger.Warn().Str("filename", header.Filename).Str("ext", ext).Msg("Invalid file extension")
		h.writeError(w, http.StatusBadRequest, "invalid_file_type", "Only .xlsx files are accepted")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "spreadsheet") && !strings.Contains(contentType, "excel") && !strings.Contains(contentType, "octet-stream") {
		h.logger.Warn().Str("content_type", contentType).Msg("Invalid content type")
		h.writeError(w, http.StatusBadRequest, "invalid_content_type", "Invalid content type for XLSX file")
		return
	}

	uploadID := uuid.New().String()

	h.logger.Info().
		Str("upload_id", uploadID).
		Str("filename", header.Filename).
		Int64("size", header.Size).
		Msg("Processing file upload")


	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read file")
		h.writeError(w, http.StatusInternalServerError, "internal_error", "Failed to read uploaded file")
		return
	}

	reader := strings.NewReader(string(fileBytes))

	result, err := h.parser.Parse(ctx, reader, uploadID)
	if err != nil {
		h.logger.Error().Err(err).Str("upload_id", uploadID).Msg("Failed to parse XLSX file")

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

	response := models.UploadResponse{
		UploadID:     uploadID,
		RowsAccepted: result.RowsAccepted,
		RowsRejected: result.RowsRejected,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UploadHandler) writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
