// Package random generates pseudo-random digit and alphanumeric strings.
package random

import "math/rand"

const (
	numberBytes = "0123456789"
)

// Number returns a string of n random decimal digits.
func Number(n int) string {
	b := make([]byte, n)
	l := len(numberBytes)
	for i := range b {
		b[i] = numberBytes[rand.Intn(l)]
	}
	return string(b)
}
