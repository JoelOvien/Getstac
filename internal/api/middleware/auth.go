package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/joelovien/go-xlsx-api/internal/models"
)

// APIKeyAuth creates a middleware for API key authentication
func APIKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If no API key is configured, skip authentication
			if apiKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check for API key in header
			providedKey := r.Header.Get("X-API-Key")
			if providedKey == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(models.ErrorResponse{
					Code:    "missing_api_key",
					Message: "API key is required",
				})
				return
			}

			// Validate API key
			if providedKey != apiKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(models.ErrorResponse{
					Code:    "invalid_api_key",
					Message: "Invalid API key",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
