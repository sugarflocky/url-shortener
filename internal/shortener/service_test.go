package shortener

import (
	"context"
	"errors"
	"testing"

	"github.com/sugarflocky/url-shortener/internal/storage"
	"github.com/sugarflocky/url-shortener/internal/storage/memory"
)

type testMemory struct {
	saveCalls int
}

func (m *testMemory) Save(ctx context.Context, url string, code string) error {
	m.saveCalls++
	return storage.ErrCodeTaken
}

func (m *testMemory) CodeByURL(ctx context.Context, url string) (code string, err error) {
	return "", storage.ErrNotFound
}

func (m *testMemory) URLByCode(ctx context.Context, code string) (url string, err error) {
	return "", storage.ErrNotFound
}

func TestCodeLength(t *testing.T) {
	ctx := context.Background()
	s := New(memory.New())

	url := "correctURL"

	code, err := s.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("Shorten(%q): %v", url, err)
	}

	got, err := s.Resolve(ctx, code)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", code, err)
	}
	if len(code) != codeLength {
		t.Errorf("Code length got %d, want %d", len(code), codeLength)
	}
	if got != url {
		t.Errorf("URL for Code %q got %q, want %q", code, got, url)
	}
}

func TestCodeIdempotency(t *testing.T) {
	ctx := context.Background()
	s := New(memory.New())

	url := "correctURL"

	code, err := s.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("Shorten(%q): %v", url, err)
	}

	got, err := s.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("Shorten(%q): %v", url, err)
	}
	if got != code {
		t.Errorf("first code %q and second code %q for same url %q are incorrect", code, got, url)
	}
}

func TestEmptyURL(t *testing.T) {
	ctx := context.Background()
	s := New(memory.New())

	url := ""

	_, err := s.Shorten(ctx, url)
	if !errors.Is(err, ErrEmptyURL) {
		t.Fatalf("Shorten for empty url: %v", err)
	}
}

func TestEmptyCode(t *testing.T) {
	ctx := context.Background()
	s := New(memory.New())

	code := ""

	_, err := s.Resolve(ctx, code)
	if !errors.Is(err, ErrEmptyCode) {
		t.Fatalf("Resolve for empty code: %v", err)
	}
}

func TestWrongCode(t *testing.T) {
	ctx := context.Background()
	s := New(memory.New())

	code := "wrongCode"

	_, err := s.Resolve(ctx, code)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("Resolve for wrong code: %v", err)
	}
}

func TestMaxAttempts(t *testing.T) {
	ctx := context.Background()
	m := &testMemory{}
	s := New(m)

	url := "correctURL"

	_, err := s.Shorten(ctx, url)
	if !errors.Is(err, ErrCodeGenerationFailed) {
		t.Fatalf("Shorten for taken code: %v", err)
	}
	if m.saveCalls != maxAttempts {
		t.Errorf("saveCalls %d and maxAttempts %d are not equal", m.saveCalls, maxAttempts)
	}
}
