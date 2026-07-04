package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sugarflocky/url-shortener/internal/shortener"
	"github.com/sugarflocky/url-shortener/internal/storage"
)

var _ shortener.Storage = (*Storage)(nil)

func TestSaveCodeAndUrl(t *testing.T) {
	ctx := context.Background()
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("url-%d", i)
			code := fmt.Sprintf("code-%d", i)

			err := s.Save(ctx, url, code)
			if err != nil {
				t.Errorf("Save(%q, %q): %v", url, code, err)
				return
			}

			got, err := s.URLByCode(ctx, code)
			if err != nil {
				t.Errorf("Got (%q by %q): %v", url, code, err)
				return
			}
			if got != url {
				t.Errorf("URLByCode(%q) = %q, want %q", code, got, url)
				return
			}
		}()
	}
	wg.Wait()
}

func TestSaveAndFind(t *testing.T) {
	ctx := context.Background()
	s := New()

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
	s := New()

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
	s := New()

	url := "correctURL"
	code := "correctCode"

	err := s.Save(ctx, url, code)
	if err != nil {
		t.Fatalf("Save(%q, %q): %v", url, code, err)
	}

	url = "URLForTakenCode"

	err = s.Save(ctx, url, code)
	if !errors.Is(err, storage.ErrCodeTaken) {
		t.Errorf("Saved for already taken code")
	}
}

func TestURLExists(t *testing.T) {
	ctx := context.Background()
	s := New()

	url := "correctURL"
	code := "correctCode"

	err := s.Save(ctx, url, code)
	if err != nil {
		t.Fatalf("Save(%q, %q): %v", url, code, err)
	}

	code = "codeForExistsURL"

	err = s.Save(ctx, url, code)
	if !errors.Is(err, storage.ErrURLExists) {
		t.Errorf("Saved for exists url")
	}
}
