package shortener

import "math/rand/v2"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const codeLength = 10

func generateShortCode() string {
	buf := make([]byte, codeLength)
	for i := range buf {
		buf[i] = alphabet[rand.IntN(len(alphabet))]
	}

	return string(buf)
}
