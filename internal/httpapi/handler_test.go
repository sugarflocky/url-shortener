package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sugarflocky/url-shortener/internal/shortener"
	"github.com/sugarflocky/url-shortener/internal/storage/memory"
)

func newTestAPI() http.Handler {
	h := New(shortener.New(memory.New()))
	return h.Router()
}

func TestValidURL(t *testing.T) {
	api := newTestAPI()

	body := strings.NewReader(`{"url":"https://testing.com/"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusCreated)
	}
	var resp codeResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response %v", err)
	}
	if len(resp.Code) != 10 {
		t.Errorf("Code length: got %d, want %d", len(resp.Code), 10)
	}
}

func TestBadJSON(t *testing.T) {
	api := newTestAPI()

	body := strings.NewReader(`{badJSON`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestInvalidURL(t *testing.T) {
	api := newTestAPI()

	body := strings.NewReader(`{"url":"incor.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestURLAndCode(t *testing.T) {
	api := newTestAPI()

	body := strings.NewReader(`{"url":"https://testing.com/"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp codeResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	rec = httptest.NewRecorder()

	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusMovedPermanently)
	}
	if loc := rec.Header().Get("Location"); loc != "https://testing.com/" {
		t.Errorf("Location got %q, want %q", loc, "https://testing.com/")
	}

}

func TestBadCode(t *testing.T) {
	api := newTestAPI()

	code := "wrongCode"
	req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
