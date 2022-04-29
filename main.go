package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	file, err := ioutil.ReadFile("response_packet.txt")
	if err != nil {
		panic("error reading file")
	}

	buf := newPacketBuffer()
	buf.buffer = file
	buf.pos = 0

	packet := NewDnsPacket()
  packet.Fill(buf)
	fmt.Printf("%+v\n", packet.Header)
}
