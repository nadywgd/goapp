package util

import (
	"math/rand"
)

var randx = rand.NewSource(42)

// RandHexString returns a random string of length n consisting of hex characters.
func RandHexString(n int) string {
	const letterBytes = "0123456789abcdef"
	const (
		letterIdxBits = 4                    // 6 bits to represent a letter index (for hex)
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, randx.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randx.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
