package main

import (
	"errors"
	"strings"
)

var (
	ErrEndOfBuffer = errors.New("end of buffer")
	ErrExceedJumps = errors.New("max amount of jumps have been exceeded")
)

// BufferSize is the length of the packet buffer
const BufferSize = 512

// Buffer represents the DNS packets buffer
type Buffer struct {
	pos int
	buf []byte
}

// NewBuffer returns a
func NewBuffer() *Buffer {
	return &Buffer{
		pos: 0,
		buf: make([]byte, BufferSize),
	}
}

// Step adds amount to the buffer position
func (b *Buffer) Step(amount int) {
	b.pos += amount
}

// Seek sets the buffer position to pos
func (b *Buffer) Seek(pos int) {
	b.pos = pos
}

// Read a single byte from the buffer and move the position forward
func (b *Buffer) Read() (byte, error) {
	var res byte
	if b.pos >= BufferSize {
		return res, ErrEndOfBuffer
	}
	res = b.buf[b.pos]
	b.pos += 1

	return res, nil
}

// Get a single byte without modifying the position property
func (b *Buffer) Get(pos int) (byte, error) {
	var res byte
	if pos >= BufferSize {
		return res, ErrEndOfBuffer
	}
	res = b.buf[pos]
	return res, nil
}

// GetRange returns bytes from a given range
func (b *Buffer) GetRange(start, length int) ([]byte, error) {
	if start+length >= BufferSize {
		return nil, ErrEndOfBuffer
	}
	return b.buf[start : start+length], nil
}

// ReadUint16 reads a uint16 valeu from 2 bytes
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

// ReadUint32 reads a uint32 value from 4 bytes
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

// ReadQName reads the buffer's qname into a domain
func (b *Buffer) ReadQName(domain *string) error {
	var (
		pos        int    = b.pos
		jumped     bool   = false
		max_jumps  int    = 5
		jump_count int    = 0
		delim      string = ""
	)

	for {
		if jump_count > max_jumps {
			return ErrExceedJumps
		}

		length, err := b.Get(pos)
		if err != nil {
			return err
		}

		if (length & 0xC0) == 0xC0 {
			if !jumped {
				b.Seek(pos + 2)
			}

			secByte, err := b.Get(pos + 1)
			if err != nil {
				return err
			}
			offset := ((uint16(length) ^ 0xC0) << 8) | uint16(secByte)
			pos = int(offset)

			jumped = true
			jump_count++
			continue
		} else {
			pos++
			if length == 0 {
				break
			}

			*domain += delim
			buffer, err := b.GetRange(pos, int(length))
			if err != nil {
				return err
			}
			*domain += strings.ToLower(string(buffer))
			delim = "."
			pos += int(length)
		}
	}

	if !jumped {
		b.Seek(pos)
	}

	return nil
}
