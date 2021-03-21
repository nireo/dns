package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

type Question struct {
	Name  string
	QType QueryType
}

func (q *Question) Read(buf *Buffer) error {
	buf.ReadQName(&q.Name)
	val, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	q.QType = GetQueryType(int(val))
	return nil
}

type Record struct {
	Domain     string
	Addr       net.IP
	TTL        uint32
	DataLength uint16
	Type       QueryType
}

func ReadRecord(buf *Buffer) (*Record, error) {
	domain := ""
	buf.ReadQName(&domain)

	qtypeVal, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}
	qtype := GetQueryType(int(qtypeVal))
	buf.ReadUint16()

	ttl, err := buf.ReadUint32()
	if err != nil {
		return nil, err
	}

	length, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}

	switch qtype {
	case A:
		addr, err := buf.ReadUint32()
		if err != nil {
			return nil, err
		}
		ipv4Addr := net.IPv4(
			byte((addr>>24)&0xFF),
			byte((addr>>16)&0xFF),
			byte((addr>>8)&0xFF),
			byte((addr>>0)&0xFF),
		)

		return &Record{
			Domain: domain,
			Addr:   ipv4Addr,
			TTL:    ttl,
			Type:   A,
		}, nil
	default:
		buf.Step(int(length))
		return &Record{
			Domain:     domain,
			Type:       qtype,
			TTL:        ttl,
			DataLength: length,
		}, nil
	}
}

type Packet struct {
	Header      *Header
	Questions   []*Question
	Answers     []*Record
	Authorities []*Record
	Resources   []*Record
}

func NewPacket() *Packet {
	return &Packet{
		Header:      NewHeader(),
		Questions:   make([]*Question, 0),
		Answers:     make([]*Record, 0),
		Authorities: make([]*Record, 0),
		Resources:   make([]*Record, 0),
	}
}

func ReadFromBuffer(buf *Buffer) (*Packet, error) {
	res := NewPacket()
	if err := res.Header.Read(buf); err != nil {
		return nil, err
	}

	for i := 0; i < int(res.Header.Questions); i++ {
		question := &Question{"", UNKNOWN}
		if err := question.Read(buf); err != nil {
			return nil, err
		}
		res.Questions = append(res.Questions, question)
	}

	for i := 0; i < int(res.Header.Answers); i++ {
		rec, err := ReadRecord(buf)
		if err != nil {
			return nil, err
		}
		res.Answers = append(res.Answers, rec)
	}

	for i := 0; i < int(res.Header.AuthoritativeEntries); i++ {
		rec, err := ReadRecord(buf)
		if err != nil {
			return nil, err
		}
		res.Authorities = append(res.Authorities, rec)
	}

	for i := 0; i < int(res.Header.ResourceEntries); i++ {
		rec, err := ReadRecord(buf)
		if err != nil {
			return nil, err
		}
		res.Resources = append(res.Resources, rec)
	}

	return res, nil
}

func main() {
	file, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatalf("error reading test file.")
	}
	buffer := &Buffer{
		pos: 0,
		buf: file,
	}

	packet, err := ReadFromBuffer(buffer)
	if err != nil {
		log.Fatalf("erro reading packet: %s", err)
	}

	for _, q := range packet.Questions {
		fmt.Printf("%v", q)
	}

	for _, r := range packet.Answers {
		fmt.Printf("%v", r)
	}

	for _, r := range packet.Authorities {
		fmt.Printf("%v", r)
	}

	for _, r := range packet.Resources {
		fmt.Printf("%v", r)
	}
}
