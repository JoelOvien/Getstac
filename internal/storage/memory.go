package storage

import (
	"sync"

	"github.com/joelovien/go-xlsx-api/internal/models"
)

// MemoryStorage provides thread-safe in-memory storage for records
type MemoryStorage struct {
	mu      sync.RWMutex
	records []models.Record
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		records: make([]models.Record, 0),
	}
}

// Store adds records to the storage
func (s *MemoryStorage) Store(records []models.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, records...)
	return nil
}

// List retrieves records with pagination
func (s *MemoryStorage) List(limit, offset int) ([]models.Record, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.records)

	// Handle edge cases
	if offset >= total {
		return []models.Record{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	// Create a copy of the slice to avoid race conditions
	result := make([]models.Record, end-offset)
	copy(result, s.records[offset:end])

	return result, total, nil
}

// Count returns the total number of records
func (s *MemoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.records)
}

// Clear removes all records (useful for testing)
func (s *MemoryStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = make([]models.Record, 0)
}

// GetByUploadID retrieves all records for a specific upload
func (s *MemoryStorage) GetByUploadID(uploadID string) []models.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Record, 0)
	for _, record := range s.records {
		if record.UploadID == uploadID {
			result = append(result, record)
		}
	}

	return result
}
