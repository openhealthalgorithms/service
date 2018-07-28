package tools

import (
	"math/rand"
	"time"
)

// charset is a set of possible symbols.
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a pseudo-random string of a given length.
func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

// StringWithCharset returns a pseudo-random string of a given length from a given charset.
func StringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
