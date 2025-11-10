package storage

import (
	"sync"

	"github.com/joelovien/go-xlsx-api/internal/models"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	records []models.Record
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		records: make([]models.Record, 0),
	}
}

func (s *MemoryStorage) Store(records []models.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, records...)
	return nil
}

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

	result := make([]models.Record, end-offset)
	copy(result, s.records[offset:end])

	return result, total, nil
}

func (s *MemoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.records)
}

func (s *MemoryStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = make([]models.Record, 0)
}

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
