package main

type ResultCode int
type QueryType int

const (
	NOERROR ResultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED

	A       QueryType = 1
	UNKNOWN QueryType = 0
)

func GetQueryType(n int) QueryType {
	if n == 1 {
		return A
	}

	return UNKNOWN
}

func GetResultCode(val uint8) ResultCode {
	switch val {
	case 1:
		return FORMERR
	case 2:
		return SERVFAIL
	case 3:
		return NXDOMAIN
	case 4:
		return NOTIMP
	case 5:
		return REFUSED
	default:
		return NOERROR
	}
}

type Header struct {
	ID                   uint16
	RecursionDesired     bool
	TruncatedMessage     bool
	AuthoritativeAnswer  bool
	Opcode               uint8
	Response             bool
	ResCode              ResultCode
	CheckingDisabled     bool
	AuthedData           bool
	Z                    bool
	RecursionAvailable   bool
	Questions            uint16
	Answers              uint16
	AuthoritativeEntries uint16
	ResourceEntries      uint16
}

func NewHeader() *Header {
	return &Header{
		ID:                   0,
		RecursionDesired:     false,
		TruncatedMessage:     false,
		AuthoritativeAnswer:  false,
		Opcode:               0,
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

func (h *Header) Read(buf *Buffer) error {
	var err error
	h.ID, err = buf.ReadUint16()
	if err != nil {
		return nil
	}

	flags, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	a := uint8(flags >> 8)
	b := uint8(flags & 0xFF)

	h.RecursionDesired = (a & (1 << 0)) > 0
	h.TruncatedMessage = (a & (1 << 1)) > 0
	h.AuthoritativeAnswer = (a & (1 << 2)) > 0
	h.Opcode = ((a >> 3) & 0x0F)
	h.Response = (a & (1 << 7)) > 0

	h.ResCode = GetResultCode(b & 0x0F)
	h.CheckingDisabled = (b & (1 << 4)) > 0
	h.AuthedData = (b & (1 << 5)) > 0
	h.Z = (b & (1 << 6)) > 0
	h.RecursionAvailable = (b & (1 << 7)) > 0

	questions, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	h.Questions = questions

	answers, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	h.Answers = answers

	authEntries, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	h.AuthoritativeEntries = authEntries

	resourceEntries, err := buf.ReadUint16()
	if err != nil {
		return err
	}
	h.ResourceEntries = resourceEntries

	return nil
}
