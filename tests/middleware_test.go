package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joelovien/go-xlsx-api/internal/api/middleware"
	"github.com/joelovien/go-xlsx-api/internal/models"
)

func TestAPIKeyAuth(t *testing.T) {
	tests := []struct {
		name           string
		configuredKey  string
		providedKey    string
		expectedStatus int
	}{
		{
			name:           "valid API key",
			configuredKey:  "secret123",
			providedKey:    "secret123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid API key",
			configuredKey:  "secret123",
			providedKey:    "wrong",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing API key",
			configuredKey:  "secret123",
			providedKey:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "no API key configured",
			configuredKey:  "",
			providedKey:    "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with auth middleware
			authMiddleware := middleware.APIKeyAuth(tt.configuredKey)
			handler := authMiddleware(testHandler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.providedKey != "" {
				req.Header.Set("X-API-Key", tt.providedKey)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		requestCount   int
		expectedStatus []int
	}{
		{
			name:         "within limit",
			requestCount: 2,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
			},
		},
		{
			name:         "exceeds limit",
			requestCount: 3,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusTooManyRequests,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh limiter for each test
			limiter := middleware.NewRateLimiter(2)
			handler := limiter.Middleware()(testHandler)

			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.RemoteAddr = "127.0.0.1:12345" // Same IP for all requests
				w := httptest.NewRecorder()

				handler.ServeHTTP(w, req)

				if w.Code != tt.expectedStatus[i] {
					t.Errorf("Request %d: expected status code %d, got %d",
						i+1, tt.expectedStatus[i], w.Code)
				}

				// Check error response for rate limited requests
				if w.Code == http.StatusTooManyRequests {
					var errResp models.ErrorResponse
					if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
						t.Fatalf("Failed to decode error response: %v", err)
					}
					if errResp.Code != "rate_limit_exceeded" {
						t.Errorf("Expected error code 'rate_limit_exceeded', got '%s'", errResp.Code)
					}
				}
			}
		})
	}
}

func TestTimeout(t *testing.T) {
	// Create a handler that takes longer than the timeout
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			// Context was cancelled
			return
		case <-time.After(200 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	})

	timeoutMiddleware := middleware.Timeout(100 * time.Millisecond)
	handler := timeoutMiddleware(slowHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// The context should be cancelled, but we can't easily check the response
	// as the handler doesn't write anything when context is cancelled
	// This test mainly ensures the middleware doesn't panic
}
