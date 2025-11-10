package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joelovien/go-xlsx-api/internal/api/handlers"
	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/joelovien/go-xlsx-api/internal/storage"
	"github.com/rs/zerolog"
)

func TestHealthHandler_Handle(t *testing.T) {
	handler := handlers.NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}
}

func TestListHandler_Handle(t *testing.T) {
	logger := zerolog.Nop()
	store := storage.NewMemoryStorage()

	// Add test data
	testRecords := []models.Record{
		{ID: "1", UploadID: "upload-1", Data: map[string]interface{}{"name": "John"}},
		{ID: "2", UploadID: "upload-1", Data: map[string]interface{}{"name": "Jane"}},
		{ID: "3", UploadID: "upload-1", Data: map[string]interface{}{"name": "Bob"}},
	}
	store.Store(testRecords)

	handler := handlers.NewListHandler(store, &logger)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
		expectedTotal  int
	}{
		{
			name:           "list all records",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
			expectedTotal:  3,
		},
		{
			name:           "list with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectedTotal:  3,
		},
		{
			name:           "list with offset",
			queryParams:    "?limit=2&offset=2",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			expectedTotal:  3,
		},
		{
			name:           "invalid limit",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid offset",
			queryParams:    "?offset=invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/records"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.ListRecordsResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if len(response.Records) != tt.expectedCount {
					t.Errorf("Expected %d records, got %d", tt.expectedCount, len(response.Records))
				}

				if response.Total != tt.expectedTotal {
					t.Errorf("Expected total %d, got %d", tt.expectedTotal, response.Total)
				}
			}
		})
	}
}
