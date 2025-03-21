package safego

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestID_Concurrent(t *testing.T) {
	const n = 100
	var wg sync.WaitGroup
	errCh := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := ID()
			if id == "" {
				errCh <- fmt.Errorf("goroutine ID should not be empty")
				return
			}
			if _, err := strconv.ParseInt(id, 10, 64); err != nil {
				// %q : 带双引号的字符串或字符字面量表示
				errCh <- fmt.Errorf("goroutine ID should be numeric, got %q", id)
			}

			t.Log(id)
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}
}
