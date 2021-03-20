package main

import "errors"

var (
	ErrEndOfBuffer = errors.New("end of buffer")
)

type Buffer struct {
	pos int
	buf []byte
}

func NewBuffer() *Buffer {
	return &Buffer{
		pos: 0,
		buf: make([]byte, 0),
	}
}

func (b *Buffer) Step(amount int) {
	b.pos += amount
}

func (b *Buffer) Seek(pos int) {
	b.pos = pos
}

// Read a single byte from the buffer and move the position forward
func (b *Buffer) Read() (byte, error) {
	var res byte
	if b.pos >= 512 {
		return res, ErrEndOfBuffer
	}
	res = b.buf[b.pos]
	b.pos += 1

	return res, nil
}

// Get a single byte without modifying the position property
func (b *Buffer) Get() (byte, error) {
	var res byte
	if b.pos >= 512 {
		return res, ErrEndOfBuffer
	}
	res = b.buf[b.pos]

	return res, nil
}

// GetRange returns bytes from a given range
func (b *Buffer) GetRange(start, length int) ([]byte, error) {
	if start+length >= 512 {
		return nil, ErrEndOfBuffer
	}
	return b.buf[start : start+length], nil
}

func (b *Buffer) ReadUint16() (uint16, error) {
	first, err := b.Read()
	if err != nil {
		return uint16(0), err
	}

	second, err := b.Read()
	if err != nil {
		return uint16(0), err
	}
	res := (uint16(first) << 8) | (uint16(second))
	return res, nil
}

func (b *Buffer) ReadUint32() (uint32, error) {
	var res uint32
	for i := 0; i < 4; i++ {
		val, err := b.Read()
		if err != nil {
			return res, err
		}
		res |= (uint32(val) >> uint32(24-i*8))
	}
	return res, nil
}
