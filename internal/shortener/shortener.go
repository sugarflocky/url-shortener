// Package shortener creates short codes for URLs
// and resolves them back to original URLs.
package shortener

import "context"

// Storage is what the shortener service needs from a storage backend.
// Implementations must be safe for concurrent use.
type Storage interface {
	// Save atomically checks and stores the url-code pair.
	// It returns storage.ErrCodeTaken if the code is taken,
	// and storage.ErrURLExists if the URL already has a code.
	Save(ctx context.Context, url string, code string) error

	// CodeByURL returns the code saved for url, or storage.ErrNotFound.
	CodeByURL(ctx context.Context, url string) (code string, err error)

	// URLByCode returns the original URL for code, or storage.ErrNotFound.
	URLByCode(ctx context.Context, code string) (url string, err error)
}
