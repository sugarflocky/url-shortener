package shortener

import "context"

type Storage interface {
	Save(ctx context.Context, url string, code string) error
	CodeByURL(ctx context.Context, url string) (code string, err error)
	URLByCode(ctx context.Context, code string) (url string, err error)
}
