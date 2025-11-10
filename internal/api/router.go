package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joelovien/go-xlsx-api/internal/api/handlers"
	custommw "github.com/joelovien/go-xlsx-api/internal/api/middleware"
	"github.com/joelovien/go-xlsx-api/internal/config"
	"github.com/joelovien/go-xlsx-api/internal/storage"
	"github.com/joelovien/go-xlsx-api/internal/xlsx"
	"github.com/rs/zerolog"
)

// NewRouter creates and configures the HTTP router
func NewRouter(cfg *config.Config, logger *zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	// Initialize storage and parser
	store := storage.NewMemoryStorage()
	parser := xlsx.NewParser(cfg.WorkerPoolSize)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(store, parser, cfg.MaxUploadSizeMB, logger)
	listHandler := handlers.NewListHandler(store, logger)
	healthHandler := handlers.NewHealthHandler()

	// Initialize rate limiter
	rateLimiter := custommw.NewRateLimiter(cfg.RateLimit)

	// Global middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(custommw.Logger(logger))
	r.Use(custommw.Timeout(cfg.RequestTimeout))

	// Health endpoint (no auth required)
	r.Get("/healthz", healthHandler.Handle)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Apply rate limiting to API routes
		r.Use(rateLimiter.Middleware())

		// Apply API key authentication if configured
		if cfg.APIKey != "" {
			r.Use(custommw.APIKeyAuth(cfg.APIKey))
		}

		// Upload endpoint
		r.Post("/uploads", uploadHandler.Handle)

		// List records endpoint
		r.Get("/records", listHandler.Handle)
	})

	return r
}
