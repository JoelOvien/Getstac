package models

import "time"

// Record represents a parsed row from the XLSX file
type Record struct {
	ID        string                 `json:"id"`
	UploadID  string                 `json:"uploadId"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"createdAt"`
}

type UploadResponse struct {
	UploadID     string `json:"uploadId"`
	RowsAccepted int    `json:"rowsAccepted"`
	RowsRejected int    `json:"rowsRejected"`
}

type ListRecordsResponse struct {
	Records []Record `json:"records"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ParsedRow struct {
	Data  map[string]interface{}
	Valid bool
	Error string
}
