package taskgo

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wangzhione/sbp/chain"
)

func DoCopyStack(a, b int) int {
	if b < 100 {
		return DoCopyStack(0, b+1)
	}
	return 0
}

func testFunc() {
	_ = DoCopyStack(0, 0)
}

var ctx = chain.Context()

const bechmarkCount = 10000000

func TestPool(t *testing.T) {
	p := NewPool(8)
	var n int32

	var wg sync.WaitGroup
	wg.Add(bechmarkCount)
	for range bechmarkCount {
		p.Go(ctx, func(context.Context) {
			defer wg.Done()
			atomic.AddInt32(&n, 1)
		})
	}
	wg.Wait()

	if n != bechmarkCount {
		t.Error(n)
	}
}

// TestPool 相对 TestGo 普通基准测试, 性能损失 1 倍, 但随着复杂业务, 二者差距没有想象那么大

func TestGo(t *testing.T) {
	var n int32

	var wg sync.WaitGroup
	wg.Add(bechmarkCount)
	for range bechmarkCount {
		go func() {
			defer wg.Done()
			atomic.AddInt32(&n, 1)
		}()
	}
	wg.Wait()

	if n != bechmarkCount {
		t.Error(n)
	}
}

func TestPoolPanic(t *testing.T) {
	testPanicFunc := func(context.Context) {
		panic("test")
	}

	p := NewPool(128)
	p.Go(ctx, testPanicFunc)

	slog.InfoContext(ctx, "Success")
}

func TestPoolWorkerExitRace(t *testing.T) {
	p := NewPool(1)
	done := make(chan struct{})

	// 模拟 worker Pop() 看到空队列准备退出，但 worker 计数尚未扣减。
	// 旧逻辑会在这个窗口里让 Go() 误判 worker 已满，导致新任务没有 worker 继续处理。
	p.worker.Store(1)
	p.Push(&task{
		ctx: ctx,
		fn: func(context.Context) {
			close(done)
		},
	})

	if !p.keepRunning() {
		t.Fatalf("fix failed: worker stopped with pending task: Len=%d, Worker=%d", p.Len(), p.Worker())
	}

	go p.running()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("fix failed: task stuck after worker exit race: Len=%d, Worker=%d", p.Len(), p.Worker())
	}
}

const benchmarkTimes = 10000

func BenchmarkPool(b *testing.B) {
	p := NewPool(int32(runtime.GOMAXPROCS(0)))

	var wg sync.WaitGroup
	b.ReportAllocs()

	for b.Loop() {
		wg.Add(benchmarkTimes)
		for range benchmarkTimes {
			p.Go(ctx, func(context.Context) {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
}

// BenchmarkPool1.895s 性能比 BenchmarkGo 1.473s

/*
goos: windows
goarch: amd64
pkg: github.com/wangzhione/sbp/helper/safego/tasks
cpu: AMD Ryzen 9 7945HX3D with Radeon Graphics

BenchmarkPool
BenchmarkPool-32
     200	   5297260 ns/op	  592224 B/op	   26981 allocs/op
PASS
ok  	github.com/wangzhione/sbp/helper/safego/tasks	1.847s

BenchmarkGo
BenchmarkGo-32
     150	   8008901 ns/op	  160483 B/op	   10001 allocs/op
PASS
ok  	github.com/wangzhione/sbp/helper/safego/tasks	2.189s

*/

func BenchmarkGo(b *testing.B) {
	var wg sync.WaitGroup
	b.ReportAllocs()

	for b.Loop() {
		wg.Add(benchmarkTimes)
		for range benchmarkTimes {
			go func() {
				testFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TestGoroutines(t *testing.T) {
	// 3 个函数分别打印 cat、dog、fish，要求每个函数都要起一个 goroutine，按照 cat、dog、fish 顺序打印在屏幕上 10 次。

	const count = 10

	catCh := make(chan struct{}, 1)
	dogCh := make(chan struct{}, 1)
	fishCh := make(chan struct{}, 1)

	var wait sync.WaitGroup
	wait.Add(3)

	fCat := func(context.Context) {
		n := 0
		for {
			n++
			t.Logf("%3d cat", n)

			dogCh <- struct{}{}

			<-catCh

			if n >= count {
				wait.Done()
				break
			}
		}
	}

	fDog := func(context.Context) {
		n := 0
		for {
			<-dogCh

			n++
			t.Logf("%3d dog", n)

			fishCh <- struct{}{}

			if n >= count {
				wait.Done()
				break
			}
		}
	}

	fFish := func(context.Context) {
		n := 0
		for {
			<-fishCh

			n++
			t.Logf("%3d fish", n)

			catCh <- struct{}{}

			if n >= count {
				wait.Done()
				break
			}
		}
	}

	p := NewPool(3)
	p.Go(ctx, fFish)
	p.Go(ctx, fDog)
	p.Go(ctx, fCat)

	wait.Wait()
}

func TestSucccessCompareInc(t *testing.T) {
	var capacity int32 = 2
	var worker int32

	var wg sync.WaitGroup
	var workerWG sync.WaitGroup
	wg.Add(20000)
	for range 20000 {
		go func() {
			defer wg.Done()

			old := atomic.LoadInt32(&worker)
			if old < capacity {
				if atomic.CompareAndSwapInt32(&worker, old, old+1) {
					current := atomic.LoadInt32(&worker)
					if current > capacity {
						t.Logf("worker=%d, capacity=%d", current, capacity)
					}

					workerWG.Add(1)
					go func() {
						defer workerWG.Done()
						defer atomic.AddInt32(&worker, -1)
					}()
				}
			}
		}()
	}
	wg.Wait()
	workerWG.Wait()
}

func TestErrorCompareInc(t *testing.T) {
	var capacity int32 = 2
	var worker int32

	var wg sync.WaitGroup
	var workerWG sync.WaitGroup
	wg.Add(400)
	for range 400 {
		go func() {
			defer wg.Done()

			if atomic.LoadInt32(&worker) < capacity {
				atomic.AddInt32(&worker, 1)

				current := atomic.LoadInt32(&worker)
				if current > capacity {
					t.Logf("worker=%d, capacity=%d", current, capacity)
				}

				workerWG.Add(1)
				go func() {
					defer workerWG.Done()
					defer atomic.AddInt32(&worker, -1)
				}()
			}
		}()
	}
	wg.Wait()
	workerWG.Wait()
}
