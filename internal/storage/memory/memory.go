// Package memory defines in-memory storage for URL-code pairs.
package memory

import (
	"context"
	"sync"

	"github.com/sugarflocky/url-shortener/internal/storage"
)

// Storage keeps URL-code pairs in memory. It is safe for concurrent use.
type Storage struct {
	mu        sync.RWMutex
	codeToURL map[string]string
	urlToCode map[string]string
}

// New creates a storage.
func New() *Storage {
	return &Storage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

// Save writes URL and code to memory.
// Returns storage.ErrCodeTaken when the code is already taken by another URL.
// Returns storage.ErrURLExists when the URL already has a code.
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

// URLByCode returns URL by given code.
// Returns storage.ErrNotFound if there is no pair.
func (s *Storage) URLByCode(ctx context.Context, code string) (url string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.codeToURL[code]
	if !ok {
		return "", storage.ErrNotFound
	}
	return url, nil
}

// CodeByURL returns code by given URL.
// Returns storage.ErrNotFound if there is no pair.
func (s *Storage) CodeByURL(ctx context.Context, url string) (code string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.urlToCode[url]
	if !ok {
		return "", storage.ErrNotFound
	}
	return code, nil
}
