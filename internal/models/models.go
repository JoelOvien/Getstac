package models

import "time"

// Record represents a parsed row from the XLSX file
type Record struct {
	ID        string                 `json:"id"`
	UploadID  string                 `json:"uploadId"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"createdAt"`
}

// UploadResponse represents the response from the upload endpoint
type UploadResponse struct {
	UploadID     string `json:"uploadId"`
	RowsAccepted int    `json:"rowsAccepted"`
	RowsRejected int    `json:"rowsRejected"`
}

// ListRecordsResponse represents the paginated response for listing records
type ListRecordsResponse struct {
	Records []Record `json:"records"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ParsedRow represents a row parsed from the XLSX file with validation status
type ParsedRow struct {
	Data  map[string]interface{}
	Valid bool
	Error string
}
