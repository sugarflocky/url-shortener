package shortener

import (
	"context"
	"errors"

	"github.com/sugarflocky/url-shortener/internal/storage"
)

var (
	ErrCodeGenerationFailed = errors.New("code generation failed")
	ErrEmptyURL             = errors.New("url is empty")
	ErrEmptyCode            = errors.New("code is empty")
)

const maxAttempts = 5 // max attempts to avoid infinity loop (basically needs only 1 attempt)

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

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

func (s *Service) Resolve(ctx context.Context, code string) (url string, err error) {
	if code == "" {
		return "", ErrEmptyCode
	}

	return s.storage.URLByCode(ctx, code)
}
