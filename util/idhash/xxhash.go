package idhash

import "github.com/cespare/xxhash/v2"

// Hash returns the hash value of the byte slice in 64bits.
func Hash(data []byte) uint64 {
	return xxhash.Sum64(data)
}

// HashString returns the hash value of the string in 64bits.
func HashString(s string) uint64 {
	// As of go 1.22, string to bytes conversion []bytes(str) is faster than using the unsafe package.
	return Hash([]byte(s))
}
