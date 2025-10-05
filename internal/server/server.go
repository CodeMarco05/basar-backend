package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/health", HealthCheck)
	// Every route her is prefixed with /api/v1 (goofy syntax mit der func, aber geil, weil man dann nicht imemr alles tippen muss)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/test", Test)
		// Add more routes
	})

	return r
}
