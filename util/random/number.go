package random

import "math/rand"

const (
	numberBytes = "0123456789"
)

// Number generates random numbers.
func Number(n int) string {
	b := make([]byte, n)
	l := len(numberBytes)
	for i := range b {
		b[i] = numberBytes[rand.Intn(l)]
	}
	return string(b)
}
