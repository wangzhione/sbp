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
	index int // range from '0' to 'len(b.data) - 1'
	begin int // beginning offset of the data stream.
	end   int // ending offset of the data stream.
}

func (b *RingBuf) Size() int  { return cap(b.data) }
func (b *RingBuf) Begin() int { return b.begin }
func (b *RingBuf) End() int   { return b.end }

func NewRingBuf(size int) RingBuf {
	return RingBuf{data: make([]byte, size)} // slice len = cap = size
}

// Reset the ring buffer
func (b *RingBuf) Reset() { b.index, b.begin, b.end = 0, 0, 0 }

// Create a copy of the buffer.
func (b *RingBuf) Dump() []byte {
	dump := make([]byte, b.Size())
	copy(dump, b.data)
	return dump
}

func (b *RingBuf) String() string {
	return fmt.Sprintf("RingBuf[size:%d, index:%d, begin:%d, end:%d]", b.Size(), b.index, b.begin, b.end)
}

// Write 环形写入 ptr []byte 数据, 优先覆盖最早的数据
// need 调用方必须保证 len(ptr) < b.Size()
func (b *RingBuf) Write(ptr []byte) {
	// 在默认 len(ptr) < b.Size() 前提下最多会有两次 copy

	// 拷贝到 b.data[b.index:] 元素个数是 min(len(b.data[b.index:]), len(ptr))
	n := copy(b.data[b.index:], ptr)
	b.end += n
	// index 索引回绕
	if b.index += n; b.index >= b.Size() {
		b.index = 0
	}
	if n < len(ptr) {
		// copy n 后续 data[n:] 到 b.data[0:]
		b.index = copy(b.data, ptr[n:])
		b.end += b.index
	}

	// 如果有效数据长度超过缓冲区容量 b.Size()，说明有旧数据被覆盖，
	// 更新 b.begin，将其移动到新的数据起点，丢弃被覆盖的旧数据。
	if b.end-b.begin > b.Size() {
		b.begin = b.end - b.Size()
	}
}

// getDataOff 逻辑上的偏移量 off 转换为缓冲区内部的实际物理偏移量 dataOff
func (b *RingBuf) getDataOff(off int) (dataOff int) {
	// off ∈ [b.begin, b.end] 要计算出相对 b.index 偏移量
	dataOff = b.index + off - b.begin
	// 偏移量超出缓冲区大小，回绕到开头
	if dataOff >= b.Size() {
		dataOff -= b.Size()
	}
	return
}

func (b *RingBuf) WriteAt(ptr []byte, off int) {
	// Low Level API 调用方必须保证 off >= b.begin && off + len(ptr) <= b.end

	writeOff := b.getDataOff(off)
	writeEnd := writeOff + b.end - off
	if writeEnd <= b.Size() {
		copy(b.data[writeOff:writeEnd], ptr)
	} else {
		n := copy(b.data[writeOff:], ptr)
		if n < len(ptr) {
			copy(b.data[:writeEnd-b.Size()], ptr[n:])
		}
	}

	// low write write 数据, 不更新 index, begin, end
}

func (b *RingBuf) EqualAt(ptr []byte, off int) bool {
	if off+len(ptr) > b.end || off < b.begin {
		return false
	}

	readOff := b.getDataOff(off)
	readEnd := readOff + len(ptr)
	if readEnd <= b.Size() {
		return bytes.Equal(ptr, b.data[readOff:readEnd])
	}

	firstLen := b.Size() - readOff
	if !bytes.Equal(ptr[:firstLen], b.data[readOff:]) {
		return false
	}
	return bytes.Equal(ptr[firstLen:], b.data[:len(ptr)-firstLen])
}

// ReadAt 与 WriteAt 互逆, read up to len(ptr), at off of the data stream.
// if return n < len(ptr) => io.EOF
func (b *RingBuf) ReadAt(ptr []byte, off int) (n int) {
	// need off ∈ [b.begin, b.end]

	readOff := b.getDataOff(off)
	readEnd := readOff + b.end - off
	if readEnd <= b.Size() {
		n = copy(ptr, b.data[readOff:readEnd])
	} else {
		n = copy(ptr, b.data[readOff:])
		if n < len(ptr) {
			n += copy(ptr[n:], b.data[:readEnd-b.Size()])
		}
	}

	return
}

// Slice returns a slice of the supplied range of the ring buffer. It will
// not alloc unless the requested range wraps the ring buffer.
func (b *RingBuf) Slice(length, off int) ([]byte, error) {
	if off > b.end || off < b.begin {
		return nil, ErrOutOfRange
	}

	readOff := b.getDataOff(off)
	readEnd := readOff + length
	if readEnd <= b.Size() {
		// slice[low:high:max] 内存共享 结果的长度为 high - low ; 结果的容量为 max - low
		// low: 切片的起始索引（包含）; high: 切片的结束索引（不包含）; max: 切片的容量上限（不包含）
		return b.data[readOff:readEnd:readEnd], nil
	}

	buf := make([]byte, length)
	n := copy(buf, b.data[readOff:])
	if n < length {
		n += copy(buf[n:], b.data[:readEnd-b.Size()])
	}
	if n < length {
		return nil, io.EOF
	}
	return buf, nil
}

// Evacuate read the data at off, then write it to the the data stream,
// Keep it from being overwritten by new data.
func (b *RingBuf) Evacuate(length int, off int) (newOff int) {
	if off < b.begin || off+length > b.end {
		return -1
	}

	readOff := b.getDataOff(off)
	if readOff == b.index {
		// no copy evacuate
		b.index += length
		if b.index >= b.Size() {
			b.index -= b.Size()
		}
	} else if readOff < b.index {
		// 情况 2: 读取位置在写入位置之前
		var n = copy(b.data[b.index:], b.data[readOff:readOff+length])
		b.index += n
		if b.index == b.Size() {
			b.index = copy(b.data, b.data[readOff+n:readOff+length])
		}
	} else {
		// 情况 3: 读取位置在写入位置之后（readOff > b.index）
		var readEnd = readOff + length
		if readEnd <= b.Size() {
			b.index += copy(b.data[b.index:], b.data[readOff:readEnd])
		} else {
			var n = copy(b.data[b.index:], b.data[readOff:])
			var tail = length - n
			b.index += n
			n = copy(b.data[b.index:], b.data[:tail])
			b.index += n
			if b.index == b.Size() {
				b.index = copy(b.data, b.data[n:tail])
			}
		}
	}

	newOff = b.end
	b.end += length
	if b.begin < b.end-b.Size() {
		b.begin = b.end - b.Size()
	}
	return
}

func (b *RingBuf) Resize(newSize int) {
	if b.Size() == newSize {
		return
	}

	newData := make([]byte, newSize) // len = cap = new size
	var offset int
	if b.end-b.begin == b.Size() {
		offset = b.index
	}
	if b.end-b.begin > newSize {
		discard := b.end - b.begin - newSize
		offset = (offset + discard) % b.Size()
		b.begin = b.end - newSize
	}

	// 如果 newSize < b.Size()
	// 只会顺序保留 [ b.index, b.Size() ) + [0, b.index) 中 newSize 个部分数据
	n := copy(newData, b.data[offset:])
	if n < newSize {
		copy(newData[n:], b.data[:offset])
	}
	b.data = newData
	b.index = 0
}

func (b *RingBuf) Skip(length int) {
	b.index += length
	for b.index >= b.Size() {
		b.index -= b.Size()
	}

	b.end += length
	if b.end-b.begin > b.Size() {
		b.begin = b.end - b.Size()
	}
}
