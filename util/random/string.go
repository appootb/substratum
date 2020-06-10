package random

import (
	"math/rand"
	"time"
)

const (
	letterBytes = numberBytes + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	letterIdxBits = 6                    // 6 bits to represent a letter index.
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits.
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits.
)

var (
	src = rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)
)

// Generate random string using masking with source.
func String(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
