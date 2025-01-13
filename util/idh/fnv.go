package idh

import (
	"hash/fnv"
	"math/bits"
)

// Hash returns the hash value of the byte slice in 64bits.
func Hash(data []byte) uint64 {
	fnv64a := fnv.New64a()
	fnv64a.Write(data)
	return fnv64a.Sum64()
}

// HashString returns the hash value of the string in 64bits.
func HashString(s string) uint64 {
	// As of go 1.22, string to bytes conversion []bytes(str) is faster than using the unsafe package.
	return Hash([]byte(s))
}

const (
	offset128Lower  = 0x62b821756295c58d
	offset128Higher = 0x6c62272e07bb0142
	prime128Lower   = 0x13b
	prime128Shift   = 24
)

// Hash128 returns the hash value of the byte slice in 128bits.
func Hash128(data []byte) (s [2]uint64) {
	for _, c := range data {
		s[1] ^= uint64(c)
		// Compute the multiplication
		s0, s1 := bits.Mul64(prime128Lower, s[1])
		s0 += s[1]<<prime128Shift + prime128Lower*s[0]
		// Update the values
		s[1] = s1
		s[0] = s0
	}

	return
}

// Hash128String returns the hash value of the string in 128bits.
func Hash128String(s string) [2]uint64 {
	return Hash128([]byte(s))
}
