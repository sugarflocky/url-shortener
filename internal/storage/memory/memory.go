package memory

import (
	"context"
	"sync"

	"github.com/sugarflocky/url-shortener/internal/storage"
)

type Storage struct {
	mu        sync.RWMutex
	codeToURL map[string]string
	urlToCode map[string]string
}

func New() *Storage {
	return &Storage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *Storage) Save(ctx context.Context, url string, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.codeToURL[code]
	if ok {
		return storage.ErrCodeTaken
	}
	_, ok = s.urlToCode[url]
	if ok {
		return storage.ErrURLExists
	}
	s.codeToURL[code] = url
	s.urlToCode[url] = code
	return nil
}

func (s *Storage) URLByCode(ctx context.Context, code string) (url string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.codeToURL[code]
	if !ok {
		return "", storage.ErrNotFound
	}
	return url, nil
}

func (s *Storage) CodeByURL(ctx context.Context, url string) (code string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.urlToCode[url]
	if !ok {
		return "", storage.ErrNotFound
	}
	return code, nil
}
