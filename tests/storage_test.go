package tests

import (
	"testing"
	"time"

	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/joelovien/go-xlsx-api/internal/storage"
)

func TestMemoryStorage_Store(t *testing.T) {
	tests := []struct {
		name    string
		records []models.Record
		wantErr bool
	}{
		{
			name: "store single record",
			records: []models.Record{
				{
					ID:        "1",
					UploadID:  "upload-1",
					Data:      map[string]interface{}{"name": "John"},
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "store multiple records",
			records: []models.Record{
				{ID: "1", UploadID: "upload-1", Data: map[string]interface{}{"name": "John"}, CreatedAt: time.Now()},
				{ID: "2", UploadID: "upload-1", Data: map[string]interface{}{"name": "Jane"}, CreatedAt: time.Now()},
			},
			wantErr: false,
		},
		{
			name:    "store empty slice",
			records: []models.Record{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewMemoryStorage()
			err := s.Store(tt.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify count
			if count := s.Count(); count != len(tt.records) {
				t.Errorf("Count() = %v, want %v", count, len(tt.records))
			}
		})
	}
}

func TestMemoryStorage_List(t *testing.T) {
	s := storage.NewMemoryStorage()

	// Prepare test data
	records := []models.Record{
		{ID: "1", UploadID: "upload-1", Data: map[string]interface{}{"name": "John"}, CreatedAt: time.Now()},
		{ID: "2", UploadID: "upload-1", Data: map[string]interface{}{"name": "Jane"}, CreatedAt: time.Now()},
		{ID: "3", UploadID: "upload-1", Data: map[string]interface{}{"name": "Bob"}, CreatedAt: time.Now()},
		{ID: "4", UploadID: "upload-1", Data: map[string]interface{}{"name": "Alice"}, CreatedAt: time.Now()},
		{ID: "5", UploadID: "upload-1", Data: map[string]interface{}{"name": "Charlie"}, CreatedAt: time.Now()},
	}
	s.Store(records)

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
		wantTotal int
		wantErr   bool
	}{
		{
			name:      "list first page",
			limit:     2,
			offset:    0,
			wantCount: 2,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "list second page",
			limit:     2,
			offset:    2,
			wantCount: 2,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "list last page",
			limit:     2,
			offset:    4,
			wantCount: 1,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "offset beyond total",
			limit:     2,
			offset:    10,
			wantCount: 0,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "list all",
			limit:     10,
			offset:    0,
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := s.List(tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("List() returned %v records, want %v", len(got), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("List() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestMemoryStorage_GetByUploadID(t *testing.T) {
	s := storage.NewMemoryStorage()

	// Prepare test data with different upload IDs
	records := []models.Record{
		{ID: "1", UploadID: "upload-1", Data: map[string]interface{}{"name": "John"}, CreatedAt: time.Now()},
		{ID: "2", UploadID: "upload-1", Data: map[string]interface{}{"name": "Jane"}, CreatedAt: time.Now()},
		{ID: "3", UploadID: "upload-2", Data: map[string]interface{}{"name": "Bob"}, CreatedAt: time.Now()},
		{ID: "4", UploadID: "upload-2", Data: map[string]interface{}{"name": "Alice"}, CreatedAt: time.Now()},
	}
	s.Store(records)

	tests := []struct {
		name      string
		uploadID  string
		wantCount int
	}{
		{
			name:      "get records for upload-1",
			uploadID:  "upload-1",
			wantCount: 2,
		},
		{
			name:      "get records for upload-2",
			uploadID:  "upload-2",
			wantCount: 2,
		},
		{
			name:      "get records for non-existent upload",
			uploadID:  "upload-3",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.GetByUploadID(tt.uploadID)
			if len(got) != tt.wantCount {
				t.Errorf("GetByUploadID() returned %v records, want %v", len(got), tt.wantCount)
			}
		})
	}
}
