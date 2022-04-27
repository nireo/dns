package main

import (
	"encoding/binary"
	"fmt"
)

type packetBuffer struct {
	buffer []byte
	pos uint64
}

func newPacketBuffer() *packetBuffer {
	return &packetBuffer{
		buffer: make([]byte, 1024),
		pos: 0,
	}
}

func (pb *packetBuffer) step(steps uint64) {
	pb.pos += steps
}

func (pb *packetBuffer) seek(pos uint64) {
	pb.pos = pos
}

func (pb *packetBuffer) read() byte {
	if pb.pos >= 1024 {
		// we could return an error, but it is complicates codes and makes it less clean,
		// because we need to handle the error.
		panic("end of buffer")
	}
	pb.pos += 1
	return pb.buffer[pb.pos - 1]
}

func (pb *packetBuffer) get() byte {
	if pb.pos >= 1024 {
		// we could return an error, but it is complicates codes and makes it less clean,
		// because we need to handle the error.
		panic("end of buffer")
	}
	return pb.buffer[pb.pos]
}

func (pb *packetBuffer) getRange(start, length uint64) []byte {
	if start + length >= 1024 {
		panic("end of buffer")
	}

	return pb.buffer[start:start + length]
}

func (pb *packetBuffer) readu16() uint16 {
	tempBuffer := make([]byte, 2)
	tempBuffer[0] = pb.read()
	tempBuffer[1] = pb.read()

	return binary.BigEndian.Uint16(tempBuffer)
}

func (pb *packetBuffer) readu32() uint32 {
	tempBuffer := make([]byte, 4)
	tempBuffer[0] = pb.read()
	tempBuffer[1] = pb.read()
	tempBuffer[2] = pb.read()
	tempBuffer[3] = pb.read()

	return binary.BigEndian.Uint32(tempBuffer)
}

func main() {
	fmt.Println("hello dns")
}
