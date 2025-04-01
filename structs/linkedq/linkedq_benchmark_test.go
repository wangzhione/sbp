package linkedq

import "testing"

// Linux + GCC
// go test -benchmem -race -run=^$ -bench ^BenchmarkLinkedQueue_EnqueueDequeue$ github.com/wangzhione/sbp/structs/linkedq -v -count=1

func BenchmarkLinkedQueue_EnqueueDequeue(b *testing.B) {
	q := New[int]()
	for i := 0; b.Loop(); i++ {
		q.Push(i)
		q.Pop()
	}
}

/*

goos: windows
goarch: amd64
pkg: github.com/wangzhione/sbp/structs/linkedq
cpu: AMD Ryzen 9 7945HX3D with Radeon Graphics
BenchmarkLinkedQueue_EnqueueDequeue
BenchmarkLinkedQueue_EnqueueDequeue-32
53748028	        21.76 ns/op	      16 B/op	       1 allocs/op

BenchmarkChannel_SendRecv
BenchmarkChannel_SendRecv-32
66517740	        17.33 ns/op	       0 B/op	       0 allocs/op

*/

func BenchmarkChannel_SendRecv(b *testing.B) {
	ch := make(chan int, 1024)
	for i := 0; b.Loop(); i++ {
		ch <- i
		<-ch
	}
}
