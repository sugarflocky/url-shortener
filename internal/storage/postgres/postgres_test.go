package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sugarflocky/url-shortener/internal/shortener"
	"github.com/sugarflocky/url-shortener/internal/storage"
)

var _ shortener.Storage = (*Storage)(nil)

// newTestStorage connects to the database from TEST_POSTGRES_DSN
// or skips the test when the variable is not set.
func newTestStorage(t *testing.T) *Storage {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN is not set")
	}

	s, err := New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	t.Cleanup(func() {
		if _, err := s.pool.Exec(context.Background(), "TRUNCATE links"); err != nil {
			t.Errorf("truncate links: %v", err)
		}
		s.pool.Close()
	})
	return s
}

func TestSaveAndFind(t *testing.T) {
	ctx := context.Background()
	s := newTestStorage(t)

	url := "correctUrl"
	code := "correctCode"

	err := s.Save(ctx, url, code)
	if err != nil {
		t.Fatalf("Save(%q, %q): %v", url, code, err)
	}

	got, err := s.CodeByURL(ctx, url)
	if err != nil {
		t.Fatalf("CodeByURL(%q): %v", url, err)
	}
	if got != code {
		t.Errorf("CodeByURL(%q) = %q, want %q", url, got, code)
	}

	got, err = s.URLByCode(ctx, code)
	if err != nil {
		t.Fatalf("URLByCode(%q): %v", code, err)
	}
	if got != url {
		t.Errorf("URLByCode(%q) = %q, want %q", code, got, url)
	}
}

func TestNotFound(t *testing.T) {
	ctx := context.Background()
	s := newTestStorage(t)

	url := "WrongURL"
	code := "WrongCode"

	_, err := s.URLByCode(ctx, code)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("URLByCode(%q): %v", code, err)
	}

	_, err = s.CodeByURL(ctx, url)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("CodeByURL(%q): %v", url, err)
	}
}

func TestCodeTaken(t *testing.T) {
	ctx := context.Background()
	s := newTestStorage(t)

	url := "correctURL"
	code := "correctCode"

	err := s.Save(ctx, url, code)
	if err != nil {
		t.Fatalf("Save(%q, %q): %v", url, code, err)
	}

	url = "URLForTakenCode"

	err = s.Save(ctx, url, code)
	if !errors.Is(err, storage.ErrCodeTaken) {
		t.Errorf("Save(%q, %q) = %v, want storage.ErrCodeTaken", url, code, err)
	}
}

func TestURLExists(t *testing.T) {
	ctx := context.Background()
	s := newTestStorage(t)

	url := "correctURL"
	code := "correctCode"

	err := s.Save(ctx, url, code)
	if err != nil {
		t.Fatalf("Save(%q, %q): %v", url, code, err)
	}

	code = "codeForExistsURL"

	err = s.Save(ctx, url, code)
	if !errors.Is(err, storage.ErrURLExists) {
		t.Errorf("Save(%q, %q) = %v, want storage.ErrURLExists", url, code, err)
	}
}
