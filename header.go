package main

type ResultCode int

const (
	NOERROR ResultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)

func GetResultCode(val int) ResultCode {
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

	// TODO: implement rest of the fields
	return nil
}
