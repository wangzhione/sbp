package chain

import (
	"sync"
	"testing"
)

func TestUUID(t *testing.T) {
	id := UUID()

	t.Logf("id = %s", id)
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
