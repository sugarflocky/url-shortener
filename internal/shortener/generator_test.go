package shortener

import (
	"strings"
	"testing"
)

func TestGenerateShortCode_Length(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := generateShortCode()

		if len(code) != codeLength {
			t.Fatalf("got code length %d, want %d", len(code), codeLength)
		}
	}
}

func TestGenerateShortCode_Alphabet(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := generateShortCode()

		for j := 0; j < len(code); j++ {
			if strings.IndexByte(alphabet, code[j]) < 0 {
				t.Fatalf("got invalid character %q in code %q", code[j], code)
			}
		}
	}
}

func TestGenerateShortCode_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		code := generateShortCode()
		if codes[code] {
			t.Fatalf("got duplicate code %q", code)
		}
		codes[code] = true
	}
}
