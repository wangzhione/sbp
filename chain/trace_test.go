package chain

import (
	"sync"
	"testing"
)

func TestCopyTrace(t *testing.T) {
	if any("X-Request-Id") == any(XRquestID) {
		t.Log("equal") // any("X-Request-Id") == any("X-Request-Id") | type equal , value equal
	} else {
		t.Log("no equal")
	}
}

func TestUUID(t *testing.T) {
	id := UUID()
	t.Logf("id = %s", id) // id = 22ba3cffc8de4a2d9dc8a95d09ed03e1

	// import "github.com/google/uuid"
	// go mod tidy
	//
	// id := uuid.New().String()
	// t.Logf("id = %s", id) // id = 22ba3cff-c8de-4a2d-9dc8-a95d09ed03e1
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
