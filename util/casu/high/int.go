// Package high provides high-performance integer parsing functions.
package high

import (
	"errors"
	"strconv"
)

// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Atoi(s []byte) (int, error) {
	sLen := len(s)
	if strconv.IntSize == 32 && (0 < sLen && sLen < 10) ||
		strconv.IntSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.

		i := 0
		switch s[0] {
		case '-', '+':
			i = 1
		}

		n := 0
		for ; i < sLen; i++ {
			ch := s[i] - '0'
			// type byte = uint8
			if ch > 9 {
				return 0, strconv.ErrSyntax
			}
			n = n*10 + int(ch)
		}
		if s[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid, big, or underscored integers.
	i64, err := ParseInt(s, strconv.IntSize)
	return int(i64), err
}

// ParseInt interprets a string s in the given base 10 and
// bit size (strconv.IntSize or 32 or 64) and returns the corresponding value i.
//
// The string may begin with a leading sign: "+" or "-".
/*
	if bitSize == 0 {
		bitSize = strconv.IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, ErrBitSize
	}

	u, err := ParseUint([]byte(s), bitSize)
*/
func ParseInt(s []byte, bitSize int) (int64, error) {
	if len(s) == 0 {
		return 0, nil
	}

	// Pick off leading sign.
	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		s = s[1:]
		neg = true
	}

	// Convert unsigned and check range.
	un, err := ParseUint(s, bitSize)
	if err != nil {
		return 0, err
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return 0, strconv.ErrRange
	}
	if neg && un > cutoff {
		return 0, strconv.ErrRange
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// ErrBitSize bit size error , need bit size in [0, 64]
var ErrBitSize = errors.New("invalid bit size")

const MaxUint64 uint64 = 1<<64 - 1

// Cutoff is the smallest number such that cutoff*base=10 > MaxUint64.
// Use compile-time constants for common cases.
const cutoff uint64 = MaxUint64/10 + 1

// ParseUint is like [ParseInt] but for unsigned numbers.
//
// A sign prefix is not permitted. bit size (0 or 32 or 64); low api 调用需要了解内部实现
/*
 	if len(s) == 0 {
		return 0, nil
	}

	if bitSize == 0 {
		bitSize = strconv.IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, ErrBitSize
	}

	u, err := ParseUint([]byte(s), bitSize)
*/
func ParseUint(s []byte, bitSize int) (uint64, error) {
	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range s {
		d := c - '0'
		// type byte = uint8
		if d > 9 {
			return 0, strconv.ErrSyntax
		}

		if n >= cutoff {
			// n*base overflows
			return 0, strconv.ErrRange
		}
		n *= 10

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return 0, strconv.ErrRange
		}
		n = n1
	}

	return n, nil
}
