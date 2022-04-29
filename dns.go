package main

type resultCode uint16
type queryType uint16

const (
	NOERROR resultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)

const (
	UNKNOWN queryType = iota
	A
)

type DnsHeader struct {
	ID                   uint16
	RecursionDesired     bool
	TruncatedMessage     bool
	AuthoritativeAnswer  bool
	OpCode               uint8
	Response             bool
	ResCode              resultCode
	CheckingDisabled     bool
	AuthedData           bool
	Z                    bool
	RecursionAvailable   bool
	Questions            uint16
	Answers              uint16
	AuthoritativeEntries uint16
	ResourceEntries      uint16
}

type DnsQuestion struct {
	Name  string
	QType queryType
}

type IPv4 struct {
	a byte
	b byte
	c byte
	d byte
}

type DnsRecord struct {
	TTL     uint32
	Domain  string
	DataLen uint16
	QType   queryType
	Addr    *IPv4 // can be null
}

type DnsPacket struct {
	Header      *DnsHeader
	Questions   []*DnsQuestion
	Answers     []*DnsRecord
	Authorities []*DnsRecord
	Resources   []*DnsRecord
}

func NewDnsQuestion(name string, qtype queryType) *DnsQuestion {
	return &DnsQuestion{
		Name:  name,
		QType: qtype,
	}
}

func NewDnsPacket() *DnsPacket {
	return &DnsPacket{
		Header:      NewDnsHeader(),
		Questions:   make([]*DnsQuestion, 0),
		Answers:     make([]*DnsRecord, 0),
		Authorities: make([]*DnsRecord, 0),
		Resources:   make([]*DnsRecord, 0),
	}
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

func NewDnsRecord() *DnsRecord {
	return &DnsRecord{
		TTL:     0,
		Domain:  "",
		DataLen: 0,
		QType:   UNKNOWN,
		Addr:    nil,
	}
}

func (dq *DnsQuestion) Fill(buffer *packetBuffer) {
	name, err := buffer.readqname()
	if err != nil {
		panic("error reading qname")
	}
	dq.Name = name
	dq.QType = queryType(buffer.readu16())
}

func (dh *DnsHeader) Fill(buffer *packetBuffer) {
	dh.ID = buffer.readu16()
	flags := buffer.readu16()
	a := uint8(flags >> 8)
	b := uint8(flags & 0xFF)
	dh.RecursionDesired = (a & (1 << 0)) > 0
	dh.TruncatedMessage = (a & (1 << 1)) > 0
	dh.AuthoritativeAnswer = (a & (1 << 2)) > 0
	dh.OpCode = (a >> 3) & 0x0F
	dh.Response = (a & (1 << 7)) > 0
	dh.ResCode = resultCode(b & 0x0F)
	dh.CheckingDisabled = (b & (1 << 4)) > 0
	dh.AuthedData = (b & (1 << 5)) > 0
	dh.Z = (b & (1 << 6)) > 0
	dh.RecursionAvailable = (b & (1 << 7)) > 0
	dh.Questions = buffer.readu16()
	dh.Answers = buffer.readu16()
	dh.AuthoritativeEntries = buffer.readu16()
	dh.ResourceEntries = buffer.readu16()
}

func (dr *DnsRecord) Fill(buffer *packetBuffer) {
	var err error
	dr.Domain, err = buffer.readqname()
	if err != nil {
		panic("error reading qname on DnsRecord")
	}

	dr.QType = queryType(buffer.readu16())
	buffer.readu16()
	dr.TTL = buffer.readu32()
	dr.DataLen = buffer.readu16()

	switch dr.QType {
	case A:
		rawAddr := buffer.readu32()
		addr := &IPv4{
			a: uint8((rawAddr >> 24) & 0xFF),
			b: uint8((rawAddr >> 16) & 0xFF),
			c: uint8((rawAddr >> 8) & 0xFF),
			d: uint8((rawAddr >> 0) & 0xFF),
		}

		dr.Addr = addr
	default:
		buffer.step(uint64(dr.DataLen))
	}
}

func (dp *DnsPacket) Fill(buffer *packetBuffer) {
	dp.Header.Fill(buffer)

	// for i := uint16(0); i < dp.Header.Questions; i += 1 {
	//	fmt.Println("reading question")
	//	question := NewDnsQuestion("", UNKNOWN)
	//	question.Fill(buffer)

	//	dp.Questions = append(dp.Questions, question)
	// }

	// for i := uint16(0); i < dp.Header.Answers; i += 1 {
	//	fmt.Println("reading answer")
	//	rec := NewDnsRecord()
	//	rec.Fill(buffer)

	//	dp.Answers = append(dp.Answers, rec)
	// }

	// for i := uint16(0); i < dp.Header.AuthoritativeEntries; i += 1 {
	//	fmt.Println("reading auth entry")
	//	rec := NewDnsRecord()
	//	rec.Fill(buffer)

	//	dp.Answers = append(dp.Authorities, rec)
	// }

	// for i := uint16(0); i < dp.Header.ResourceEntries; i += 1 {
	//	fmt.Println("reading resource")
	//	rec := NewDnsRecord()
	//	rec.Fill(buffer)

	//	dp.Answers = append(dp.Resources, rec)
	// }
}
