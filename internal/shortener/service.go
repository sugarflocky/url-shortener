package shortener

import (
	"context"
	"errors"

	"github.com/sugarflocky/url-shortener/internal/storage"
)

var (
	// ErrCodeGenerationFailed is returned by service if it cannot generate code after maxAttempts attempts.
	ErrCodeGenerationFailed = errors.New("code generation failed")
	// ErrEmptyURL is returned by service if given URL is empty.
	ErrEmptyURL = errors.New("url is empty")
	// ErrEmptyCode is returned by service if given code is empty.
	ErrEmptyCode = errors.New("code is empty")
)

const maxAttempts = 5 // max attempts to avoid infinity loop (basically needs only 1 attempt)

type Service struct {
	storage Storage
}

// New creates a service with given storage.
func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Shorten generates code for URL, returns the same code for the same URL.
// Returns ErrEmptyURL if URL is empty.
// Returns ErrCodeGenerationFailed if code cannot be created after maxAttempts attempts.
func (s *Service) Shorten(ctx context.Context, url string) (code string, err error) {
	if url == "" {
		return "", ErrEmptyURL
	}

	code, err = s.storage.CodeByURL(ctx, url)
	if err == nil {
		return code, nil
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return "", err
	}

	for i := 0; i < maxAttempts; i++ {
		code = generateShortCode()
		err = s.storage.Save(ctx, url, code)
		if errors.Is(err, storage.ErrCodeTaken) {
			continue
		}
		if errors.Is(err, storage.ErrURLExists) {
			return s.storage.CodeByURL(ctx, url)
		}

		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", ErrCodeGenerationFailed
}

// Resolve returns URL for given code.
// Returns ErrEmptyCode if code is empty.
func (s *Service) Resolve(ctx context.Context, code string) (url string, err error) {
	if code == "" {
		return "", ErrEmptyCode
	}

	return s.storage.URLByCode(ctx, code)
}
