package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	mux    *http.ServeMux
}

// NewServer creates a new HTTP server
func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		mux: mux,
	}
}

// WithMiddleware sets middleware for the server
func (s *Server) WithMiddleware(middleware func(http.Handler) http.Handler) {
    s.server.Handler = middleware(s.mux)
}

// Run starts the HTTP server
func (s *Server) Run() error {
	slog.Info("starting the API server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down the API server")
	return s.server.Shutdown(ctx)
}

// Mux returns the HTTP router
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}
