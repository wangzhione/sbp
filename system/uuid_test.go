package system

import (
	"encoding/hex"
	"sync"
	"testing"
	"time"
)

func TestUUID(t *testing.T) {
	start := time.Now().UnixMilli()
	id := UUID()
	end := time.Now().UnixMilli()

	raw := mustDecodeUUIDHex(t, id)

	if raw[6]>>4 != 0x7 {
		t.Fatalf("UUID() version = %x, want 7, id = %s", raw[6]>>4, id)
	}
	if raw[8]>>6 != 0x2 {
		t.Fatalf("UUID() variant = %b, want 10, id = %s", raw[8]>>6, id)
	}

	gotMS := int64(raw[0])<<40 | int64(raw[1])<<32 | int64(raw[2])<<24 |
		int64(raw[3])<<16 | int64(raw[4])<<8 | int64(raw[5])
	if gotMS < start || gotMS > end {
		t.Fatalf("UUID() timestamp = %d, want between %d and %d, id = %s", gotMS, start, end, id)
	}
}

func TestUUIDUnique(t *testing.T) {
	seen := make(map[string]struct{}, 10000)
	for range 10000 {
		id := UUID()
		mustDecodeUUIDHex(t, id)

		if _, ok := seen[id]; ok {
			t.Fatalf("UUID() duplicate id = %s", id)
		}
		seen[id] = struct{}{}
	}
}

func TestUUIDv4(t *testing.T) {
	id := UUIDv4()
	raw := mustDecodeUUIDHex(t, id)

	if raw[6]>>4 != 0x4 {
		t.Fatalf("UUIDv4() version = %x, want 4, id = %s", raw[6]>>4, id)
	}
	if raw[8]>>6 != 0x2 {
		t.Fatalf("UUIDv4() variant = %b, want 10, id = %s", raw[8]>>6, id)
	}
}

func mustDecodeUUIDHex(t *testing.T, id string) []byte {
	t.Helper()

	if len(id) != 32 {
		t.Fatalf("uuid length = %d, want 32, id = %s", len(id), id)
	}
	for _, c := range id {
		// UUID() 返回无横线的小写 hex，避免大小写或格式误改。
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Fatalf("uuid contains non-lower-hex char %q, id = %s", c, id)
		}
	}

	raw, err := hex.DecodeString(id)
	if err != nil {
		t.Fatalf("decode uuid hex: %v, id = %s", err, id)
	}
	if len(raw) != 16 {
		t.Fatalf("decoded uuid bytes = %d, want 16, id = %s", len(raw), id)
	}
	return raw
}

func BenchmarkUUID(b *testing.B) {
	for b.Loop() {
		_ = UUID()
	}
}

// 测试高并发下的性能表现
func BenchmarkUUIDHighConcurrency(b *testing.B) {
	b.Run("HighConcurrency_1000", func(b *testing.B) {
		var wg sync.WaitGroup
		concurrency := 1000

		b.ResetTimer()

		for range concurrency {
			wg.Go(func() {
				iterations := b.N / concurrency
				if iterations == 0 {
					iterations = 1
				}
				for j := 0; j < iterations; j++ {
					_ = UUID()
				}
			})
		}

		wg.Wait()
	})
}
