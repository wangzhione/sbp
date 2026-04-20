package casu

const digits = "0123456789"

const nSmalls = 100

// smalls is the formatting of 00..99 concatenated.
// It is then padded out with 56 x's to 256 bytes,
// so that smalls[x&0xFF] has no bounds check.
const smalls = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

const host64bit = ^uint(0)>>32 != 0

// formatBase10 formats the decimal representation of u into the tail of a
// and returns the offset of the first byte written to a. That is, after
//
//	i := formatBase10(a, u)
//
// the decimal representation is in a[i:].
func formatBase10(a []byte, u uint64) int {
	// Split into 9-digit chunks that fit in uint32s
	// and convert each chunk using uint32 math instead of uint64 math.
	// The obvious way to write the outer loop is "for u >= 1e9", but most numbers are small,
	// so the setup for the comparison u >= 1e9 is usually pure overhead.
	// Instead, we approximate it by u>>29 != 0, which is usually faster and good enough.
	i := len(a)
	for (host64bit && u>>29 != 0) || (!host64bit && uint32(u)>>29|uint32(u>>32) != 0) {
		var lo uint32
		u, lo = u/1e9, uint32(u%1e9)

		// Convert 9 digits.
		for range 4 {
			var dd uint32
			lo, dd = lo/100, (lo%100)*2
			i -= 2
			a[i+0], a[i+1] = smalls[dd+0], smalls[dd+1]
		}
		i--
		a[i] = smalls[lo*2+1]

		// If we'd been using u >= 1e9 then we would be guaranteed that u/1e9 > 0,
		// but since we used u>>29 != 0, u/1e9 might be 0, so we might be done.
		// (If u is now 0, then at the start we had 2²⁹ ≤ u < 10⁹, so it was still correct
		// to write 9 digits; we have not accidentally written any leading zeros.)
		if u == 0 {
			return i
		}
	}

	// Convert final chunk, at most 8 digits.
	lo := uint32(u)
	for lo >= 100 {
		var dd uint32
		lo, dd = lo/100, (lo%100)*2
		i -= 2
		a[i+0], a[i+1] = smalls[dd+0], smalls[dd+1]
	}
	i--
	dd := lo * 2
	a[i] = smalls[dd+1]
	if lo >= 10 {
		i--
		a[i] = smalls[dd+0]
	}
	return i
}

func format10(u uint64, neg bool) string {
	// 0 ~ 2^64 - 1 = 18,446,744,073,709,551,615
	// + - for sign of 64 bit value in base 10, max len 20 + 1 ;
	var a [20 + 1]byte

	if neg {
		u = -u
	}

	i := formatBase10(a[:], u)

	if neg {
		i--
		a[i] = '-'
	}

	return string(a[i:])
}
