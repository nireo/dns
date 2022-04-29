package main

import "bytes"

type resultCode int

const (
	NOERROR resultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)

type DnsHeader struct {
	ID uint16

	RecursionDesired    bool
	TruncatedMessage    bool
	AuthoritativeAnswer bool
	OpCode              uint8
	Response            bool

	ResCode            resultCode
	CheckingDisabled   bool
	AuthedData         bool
	Z                  bool
	RecursionAvailable bool

	Questions            uint16
	Answers              uint16
	AuthoritativeEntries uint16
	ResourceEntries      uint16
}

func NewDnsHeader() *DnsHeader {
	return &DnsHeader{
		ID:                   0,
		RecursionDesired:     false,
		TruncatedMessage:     false,
		AuthoritativeAnswer:  false,
		OpCode:               0,
		Response:             false,
		ResCode:              NOERROR,
		CheckingDisabled:     false,
		AuthedData:           false,
		Z:                    false,
		RecursionAvailable:   false,
		Questions:            0,
		Answers:              0,
		AuthoritativeEntries: 0,
		ResourceEntries:      0,
	}
}

func (dh *DnsHeader) FillHeader(buffer *packetBuffer) {
}
