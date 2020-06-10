package random

const (
	numberBytes = "0123456789"
)

// Generate random numbers.
func Number(n int) string {
	b := make([]byte, n)
	l := len(numberBytes)
	for i := range b {
		b[i] = numberBytes[rnd.Intn(l)]
	}
	return string(b)
}
