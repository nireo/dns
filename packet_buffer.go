package main

import (
	"encoding/binary"
	"fmt"
)

type packetBuffer struct {
	buffer []byte
	pos    uint64
}

func newPacketBuffer() *packetBuffer {
	return &packetBuffer{
		buffer: make([]byte, 1024),
		pos:    0,
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
	return pb.buffer[pb.pos-1]
}

func (pb *packetBuffer) get(pos uint64) byte {
	if pos >= uint64(len(pb.buffer)) {
		// we could return an error, but it is complicates codes and makes it less clean,
		// because we need to handle the error.
		panic("end of buffer")
	}
	return pb.buffer[pos]
}

func (pb *packetBuffer) getRange(start, length uint64) []byte {
	if start+length >= 1024 {
		panic("end of buffer")
	}

	return pb.buffer[start : start+length]
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

func (pb *packetBuffer) readqname() (string, error) {
	pos := pb.pos

	jumped := false
	max_jumps := 5
	jumps_done := 0

	outStr := ""
	delim := ""
	for {
		if jumps_done > max_jumps {
			return "", fmt.Errorf("Limit of %d jumps exceeded.", max_jumps)
		}

		len := pb.get(pos)
		if (len & 0xC0) == 0xC0 {
			if !jumped {
				pb.seek(pos + 2)
			}

			byte2 := uint16(pb.get(pos + 1))
			offset := ((uint16(len) ^ 0xC0) << 8) | byte2
			pos = uint64(offset)

			jumped = true
			jumps_done += 1

			continue
		} else {
			pos += 1
			if len == 0 {
				break
			}

			outStr += delim
			strBuffer := pb.getRange(pos, uint64(len))
			outStr += string(strBuffer)
			delim = "."

			pos += uint64(len)
		}
	}

	if !jumped {
		pb.seek(pos)
	}

	return outStr, nil
}
