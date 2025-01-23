package cache

import (
	"testing"
)

func TestRingBuf(t *testing.T) {
	rb := NewRingBuf(16)
	for i := 0; i < 2; i++ {
		rb.Write([]byte("fghibbbbccccddde"))
		rb.Write([]byte("fghibbbbc"))
		rb.Resize(16)
		off := rb.Evacuate(3, 9)
		t.Log(string(rb.Dump()))
		if off != rb.End()-3 {
			t.Log(string(rb.Dump()), rb.End())
			t.Fatalf("off got %v", off)
		}
		off = rb.Evacuate(5, 15)
		t.Log(string(rb.Dump()))
		if off != rb.End()-5 {
			t.Fatalf("off got %v", off)
		}
		rb.Resize(64)
		rb.Resize(32)
		data := make([]byte, 5)
		rb.ReadAt(data, off)
		if string(data) != "efghi" {
			t.Fatalf("read at should be efghi, got %v", string(data))
		}

		off = rb.Evacuate(10, 0)
		if off != -1 {
			t.Fatal("evacutate out of range offset should return error")
		}

		/* -- After reset the buffer should behave exactly the same as a new one.
		 *    Hence, run the test once more again with reset buffer. */
		rb.Reset()
	}

	/*
		true
		   fghibbbbccccddde
		   ighibbbbccccefgh
		   cccdddefghibbbbc
		   cccefghighibbbbc

	*/
}

func TestRingBuf_Dump(t *testing.T) {
	rb := NewRingBuf(16)

	t.Log("rb.Size()", rb.Size())

	data := rb.Dump()
	t.Log("len(data)", len(data))

	if rb.Size() != len(data) {
		t.Error("rb.Size() != len(data)")
	}
}

func TestRingBuf_Write(t *testing.T) {
	rb := NewRingBuf(5)

	t.Log(rb.String(), rb.data)

	rb.Write([]byte("123"))
	t.Log(rb.String(), rb.data)

	rb.Write([]byte("4567"))
	t.Log(rb.String(), rb.data)

	rb.Write([]byte("89"))
	t.Log(rb.String(), rb.data)

	rb.Write([]byte("0"))
	t.Log(rb.String(), rb.data)

	rb.Write([]byte("1"))
	t.Log(rb.String(), rb.data)

	/*
	   RingBuf[size:5, index:0, begin:0, end:0] [0 0 0 0 0]
	   RingBuf[size:5, index:3, begin:0, end:3] [49 50 51 0 0]
	   RingBuf[size:5, index:2, begin:2, end:7] [54 55 51 52 53]
	   RingBuf[size:5, index:4, begin:4, end:9] [54 55 56 57 53]
	   RingBuf[size:5, index:0, begin:5, end:10] [54 55 56 57 48]
	   RingBuf[size:5, index:1, begin:6, end:11] [49 55 56 57 48]
	*/
}

func TestRingBuf_WriteLimit(t *testing.T) {
	var rbSize int16 = 5

	var rbBegin, rbEnd int16

	t.Log("rbBegin", rbBegin, "rbEnd", rbEnd, "rbEnd-rbBegin", rbEnd-rbBegin)

	for {
		rbEnd += 3
		if rbEnd-rbBegin > rbSize {
			rbBegin = rbEnd - rbSize
		}

		if rbEnd < 0 {
			t.Log("rbBegin", rbBegin, "rbEnd", rbEnd, "rbEnd-rbBegin", rbEnd-rbBegin)
			break
		}
	}

	rbEnd += 3
	if rbEnd-rbBegin > rbSize {
		rbBegin = rbEnd - rbSize
	}
	t.Log("rbBegin", rbBegin, "rbEnd", rbEnd, "rbEnd-rbBegin", rbEnd-rbBegin)

	rbEnd += 3
	if rbEnd-rbBegin > rbSize {
		rbBegin = rbEnd - rbSize
	}
	t.Log("rbBegin", rbBegin, "rbEnd", rbEnd, "rbEnd-rbBegin", rbEnd-rbBegin)

	/*
	   rbBegin 0 rbEnd 0 rbEnd-rbBegin 0
	   rbBegin 32764 rbEnd -32767 rbEnd-rbBegin 5
	   rbBegin 32767 rbEnd -32764 rbEnd-rbBegin 5
	   rbBegin -32766 rbEnd -32761 rbEnd-rbBegin 5
	*/
}
