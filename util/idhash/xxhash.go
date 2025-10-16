package idhash

import "github.com/cespare/xxhash/v2"

// Hash returns the hash value of the byte slice in 64bits.
// Hash String returns the hash value of the string in 64bits.
func Hash[T ~[]byte | ~string](data T) uint64 {
	return xxhash.Sum64([]byte(data))
}
