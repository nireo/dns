package main

import (
	"fmt"
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

func main() {
	fmt.Println("running dns server...")
	select {}
}
