package cache

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

/*
 * RingBuf 这是个 Low Level 级别库, 使用起来需要用的人有阅读和理解源代码
 */

// ErrOutOfRange out of range 索引越界
var ErrOutOfRange = errors.New("out of range")

// Ring buffer has a fixed size, when data exceeds the
// size, old data will be overwritten by new data.
// It only contains the data in the stream from begin to end
type RingBuf struct {
	data  []byte
	index int // range from '0' to 'cap(rb.data) - 1'
	begin int // beginning offset of the data stream.
	end   int // ending offset of the data stream.
}

func (rb *RingBuf) Size() int {
	return cap(rb.data)
}

func (rb *RingBuf) Begin() int {
	return rb.begin
}

func (rb *RingBuf) End() int {
	return rb.end
}

func NewRingBuf(size int) RingBuf {
	return RingBuf{data: make([]byte, size)} // slice len = cap = size
}

// Reset the ring buffer
func (rb *RingBuf) Reset() {
	rb.index = 0
	rb.begin, rb.end = 0, 0
}

// Create a copy of the buffer.
func (rb *RingBuf) Dump() []byte {
	dump := make([]byte, rb.Size())
	copy(dump, rb.data)
	return dump
}

func (rb *RingBuf) String() string {
	return fmt.Sprintf("RingBuf[size:%d, index:%d, begin:%d, end:%d]", rb.Size(), rb.index, rb.begin, rb.end)
}

// Write 环形写入 ptr []byte 数据, 优先覆盖最早的数据
// need 调用方必须保证 len(ptr) < rb.Size()
func (rb *RingBuf) Write(ptr []byte) {
	// 在默认 len(ptr) < rb.Size() 前提下最多会有两次 copy

	// 拷贝到 rb.data[rb.index:] 元素个数是 min(len(rb.data[rb.index:]), len(ptr))
	n := copy(rb.data[rb.index:], ptr)
	rb.end += n
	// index 索引回绕
	if rb.index += n; rb.index >= rb.Size() {
		rb.index = 0
	}
	if n < len(ptr) {
		// copy n 后续 data[n:] 到 rb.data[0:]
		rb.index = copy(rb.data, ptr[n:])
		rb.end += rb.index
	}

	// 如果有效数据长度超过缓冲区容量 rb.Size()，说明有旧数据被覆盖，
	// 更新 rb.begin，将其移动到新的数据起点，丢弃被覆盖的旧数据。
	if rb.end-rb.begin > rb.Size() {
		rb.begin = rb.end - rb.Size()
	}
}

// getDataOff 逻辑上的偏移量 off 转换为缓冲区内部的实际物理偏移量 dataOff
func (rb *RingBuf) getDataOff(off int) (dataOff int) {
	// off ∈ [rb.begin, rb.end] 要计算出相对 rb.index 偏移量
	dataOff = rb.index + off - rb.begin
	// 偏移量超出缓冲区大小，回绕到开头
	if dataOff >= rb.Size() {
		dataOff -= rb.Size()
	}
	return
}

func (rb *RingBuf) WriteAt(ptr []byte, off int) {
	// Low Level API 调用方必须保证 off >= rb.begin && off + len(ptr) <= rb.end

	writeOff := rb.getDataOff(off)
	writeEnd := writeOff + rb.end - off
	if writeEnd <= rb.Size() {
		copy(rb.data[writeOff:writeEnd], ptr)
	} else {
		n := copy(rb.data[writeOff:], ptr)
		if n < len(ptr) {
			copy(rb.data[:writeEnd-rb.Size()], ptr[n:])
		}
	}

	// low write write 数据, 不更新 index, begin, end
}

func (rb *RingBuf) EqualAt(p []byte, off int) bool {
	if off+len(p) > rb.end || off < rb.begin {
		return false
	}
	readOff := rb.getDataOff(off)
	readEnd := readOff + len(p)
	if readEnd <= rb.Size() {
		return bytes.Equal(p, rb.data[readOff:readEnd])
	} else {
		firstLen := rb.Size() - readOff
		equal := bytes.Equal(p[:firstLen], rb.data[readOff:])
		if equal {
			secondLen := len(p) - firstLen
			equal = bytes.Equal(p[firstLen:], rb.data[:secondLen])
		}
		return equal
	}
}

// read up to len(p), at off of the data stream.
func (rb *RingBuf) ReadAt(p []byte, off int) (n int, err error) {
	if off > rb.end || off < rb.begin {
		err = ErrOutOfRange
		return
	}

	readOff := rb.getDataOff(off)
	readEnd := readOff + rb.end - off
	if readEnd <= rb.Size() {
		n = copy(p, rb.data[readOff:readEnd])
	} else {
		n = copy(p, rb.data[readOff:])
		if n < len(p) {
			n += copy(p[n:], rb.data[:readEnd-rb.Size()])
		}
	}
	if n < len(p) {
		err = io.EOF
	}
	return
}

// Slice returns a slice of the supplied range of the ring buffer. It will
// not alloc unless the requested range wraps the ring buffer.
func (rb *RingBuf) Slice(off, length int) ([]byte, error) {
	if off > rb.end || off < rb.begin {
		return nil, ErrOutOfRange
	}
	readOff := rb.getDataOff(off)
	readEnd := readOff + length
	if readEnd <= rb.Size() {
		return rb.data[readOff:readEnd:readEnd], nil
	}
	buf := make([]byte, length)
	n := copy(buf, rb.data[readOff:])
	if n < int(length) {
		n += copy(buf[n:], rb.data[:readEnd-rb.Size()])
	}
	if n < int(length) {
		return nil, io.EOF
	}
	return buf, nil
}

// Evacuate read the data at off, then write it to the the data stream,
// Keep it from being overwritten by new data.
func (rb *RingBuf) Evacuate(off int, length int) (newOff int) {
	if off+length > rb.end || off < rb.begin {
		return -1
	}
	readOff := rb.getDataOff(off)
	if readOff == rb.index {
		// no copy evacuate
		rb.index += length
		if rb.index >= rb.Size() {
			rb.index -= rb.Size()
		}
	} else if readOff < rb.index {
		var n = copy(rb.data[rb.index:], rb.data[readOff:readOff+length])
		rb.index += n
		if rb.index == rb.Size() {
			rb.index = copy(rb.data, rb.data[readOff+n:readOff+length])
		}
	} else {
		var readEnd = readOff + length
		var n int
		if readEnd <= rb.Size() {
			n = copy(rb.data[rb.index:], rb.data[readOff:readEnd])
			rb.index += n
		} else {
			n = copy(rb.data[rb.index:], rb.data[readOff:])
			rb.index += n
			var tail = length - n
			n = copy(rb.data[rb.index:], rb.data[:tail])
			rb.index += n
			if rb.index == rb.Size() {
				rb.index = copy(rb.data, rb.data[n:tail])
			}
		}
	}
	newOff = rb.end
	rb.end += length
	if rb.begin < rb.end-rb.Size() {
		rb.begin = rb.end - rb.Size()
	}
	return
}

func (rb *RingBuf) Resize(newSize int) {
	if rb.Size() == newSize {
		return
	}
	newData := make([]byte, newSize)
	var offset int
	if rb.end-rb.begin == rb.Size() {
		offset = rb.index
	}
	if int(rb.end-rb.begin) > newSize {
		discard := int(rb.end-rb.begin) - newSize
		offset = (offset + discard) % rb.Size()
		rb.begin = rb.end - newSize
	}
	n := copy(newData, rb.data[offset:])
	if n < newSize {
		copy(newData[n:], rb.data[:offset])
	}
	rb.data = newData
	rb.index = 0
}

func (rb *RingBuf) Skip(length int) {
	rb.end += length
	rb.index += int(length)
	for rb.index >= rb.Size() {
		rb.index -= rb.Size()
	}
	if int(rb.end-rb.begin) > rb.Size() {
		rb.begin = rb.end - rb.Size()
	}
}
