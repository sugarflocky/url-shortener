package shortener

import (
	"context"
	"testing"
)

func TestCodeLength(t *testing.T) {
	ctx := context.Background()
	s := New()

	url := "correctURL"

	s.Shorten(ctx, url)

}
