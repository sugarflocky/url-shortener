package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sugarflocky/url-shortener/internal/shortener"
	"github.com/sugarflocky/url-shortener/internal/storage"
)

// Handler serves the shortener HTTP API.
type Handler struct {
	svc *shortener.Service
}

// New creates a handler with the given service.
func New(svc *shortener.Service) *Handler {
	return &Handler{svc: svc}
}

// Router returns the routes of the API.
func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.shorten)
	mux.HandleFunc("GET /{code}", h.resolve)
	return mux
}

// maxBodySize caps the request body size.
const maxBodySize = 1 << 20

type inputURLDto struct {
	URL string `json:"url"`
}

type codeResponse struct {
	Code string `json:"code"`
}

func (h *Handler) shorten(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req inputURLDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	code, err := h.svc.Shorten(r.Context(), req.URL)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(codeResponse{Code: code})
}

func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	gotURL, err := h.svc.Resolve(r.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}

	http.Redirect(w, r, gotURL, http.StatusMovedPermanently)
}

// statusFromError translates service errors into HTTP status codes.
func statusFromError(err error) int {
	switch {
	case errors.Is(err, shortener.ErrEmptyCode):
		return http.StatusBadRequest
	case errors.Is(err, shortener.ErrEmptyURL):
		return http.StatusBadRequest
	case errors.Is(err, storage.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// writeError sends the error to the client, hiding internal details.
func writeError(w http.ResponseWriter, err error) {
	status := statusFromError(err)
	msg := err.Error()
	if status == http.StatusInternalServerError {
		msg = "internal error"
	}
	http.Error(w, msg, status)
}
