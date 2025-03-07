package sets

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const N = 1000

func Test_AddConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

func Test_AppendConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)

	n := len(ints) >> 1
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			s.Append(i, N-i-1)
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

func Test_CardinalityConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		elems := s.Len()
		for i := 0; i < N; i++ {
			newElems := s.Len()
			if newElems < elems {
				t.Errorf("Cardinality shrunk from %v to %v", elems, newElems)
			}
		}
		wg.Done()
	}()

	for i := 0; i < N; i++ {
		s.Add(rand.Int())
	}
	wg.Wait()
}

func Test_ClearConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func() {
			s.Clear()
			wg.Done()
		}()
		go func(i int) {
			s.Add(i)
		}(i)
	}

	wg.Wait()
}

func Test_CloneConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)

	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := range ints {
		go func(i int) {
			s.Remove(i)
			wg.Done()
		}(i)
	}
	s.Clone()
	wg.Wait()
}

func Test_ContainssConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)
	integers := make([]int, 0)
	for _, v := range ints {
		s.Add(v)
		integers = append(integers, v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.Exists(integers...)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_ContainssOneConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	for _, v := range ints {
		number := v
		wg.Add(1)
		go func() {
			s.Contains(number)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_EqualConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewTSet[int](), NewTSet[int]()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.EQual(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_RemoveConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for _, v := range ints {
		go func(i int) {
			s.Remove(i)
			wg.Done()
		}(v)
	}
	wg.Wait()

	if s.Len() != 0 {
		t.Errorf("Expected cardinality 0; got %v", s.Len())
	}
}

func Test_StringConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for range ints {
		go func() {
			_ = s.String()
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_ToSlice(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewTSet[int]()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	setAsSlice := s.ToSlice()
	if len(setAsSlice) != s.Len() {
		t.Errorf("Set length is incorrect: %v", len(setAsSlice))
	}

	for _, i := range setAsSlice {
		if !s.Exists(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

// Test_ToSliceDeadlock - fixes issue: https://github.com/deckarep/golang-set/issues/36
// This code reveals the deadlock however it doesn't happen consistently.
func Test_ToSliceDeadlock(t *testing.T) {
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	set := NewTSet[int]()
	workers := 10
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				set.Add(1)
				set.ToSlice()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_UnmarshalJSON(t *testing.T) {
	s := []byte(`["test", "1", "2", "3"]`) //,["4,5,6"]]`)
	expected := NewTSetWithValue(
		[]string{
			string(json.Number("1")),
			string(json.Number("2")),
			string(json.Number("3")),
			"test",
		}...,
	)

	actual := NewTSet[string]()
	err := json.Unmarshal(s, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}
}

func Test_MarshalJSON(t *testing.T) {
	expected := NewTSetWithValue(
		[]string{
			string(json.Number("1")),
			"test",
		}...,
	)

	b, err := json.Marshal(
		NewTSetWithValue(
			[]string{
				"1",
				"test",
			}...,
		),
	)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	actual := NewTSet[string]()
	err = json.Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}
}
