package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/hard-gainer/url-shortener/internal/service"
	"github.com/hard-gainer/url-shortener/internal/storage"
)

// URLHandler handles URL shortening requests
type URLHandler struct {
	urlService service.URLService
	baseURL    string
}

// NewURLHandler creates a new URL handler
func NewURLHandler(urlService service.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		baseURL:    baseURL,
	}
}

// RegisterRoutes registers the handler's routes
func (h *URLHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/info/{shortURL}", h.GetURLInfo)
	mux.HandleFunc("POST /api/shorten", h.ShortenURL)
	mux.HandleFunc("GET /{shortURL}", h.HandleRequest)
}

// ShortenURLRequest is the request body for shortening a URL
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse is the response body for shortening a URL
type ShortenURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ShortenURL handles requests to create short URLs
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	if originalURL == "" {
		renderError(w, "URL is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		renderError(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.ShortenURL(r.Context(), originalURL)
	if err != nil {
		slog.Error("failed to shorten URL", "error", err)
		renderError(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	fullShortURL := h.baseURL + "/" + shortURL

	resp := ShortenURLResponse{
		ShortURL:    fullShortURL,
		OriginalURL: originalURL,
	}

	renderJSON(w, resp, http.StatusCreated)
}

// HandleRequest processes all GET requests and redirects if a valid short URL is found
func (h *URLHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	if shortURL == "" {
		renderError(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.urlService.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrURLMappingNotFound) {
			http.NotFound(w, r)
			return
		}

		slog.Error("failed to get original URL", "error", err, "short_url", shortURL)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// GetURLInfo returns information about shortened URL
func (h *URLHandler) GetURLInfo(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/api/info/")
	if shortURL == "" {
		renderError(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.urlService.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrURLMappingNotFound) {
			renderError(w, "Short URL not found", http.StatusNotFound)
			return
		}

		slog.Error("failed to get URL info", "error", err, "short_url", shortURL)
		renderError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := ShortenURLResponse{
		ShortURL:    h.baseURL + "/" + shortURL,
		OriginalURL: originalURL,
	}

	renderJSON(w, resp, http.StatusOK)
}

// renderJSON is a helper function for response formatting
func renderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// renderError is a helper function for error rendering
func renderError(w http.ResponseWriter, message string, status int) {
	resp := ErrorResponse{Error: message}
	renderJSON(w, resp, status)
}
