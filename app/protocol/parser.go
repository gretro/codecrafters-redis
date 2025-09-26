package protocol

import (
	"bufio"
	"io"
)

type RespParser struct {
	reader *bufio.Reader
}

func NewRespParser(r io.Reader) *RespParser {
	return &RespParser{
		reader: bufio.NewReader(r),
	}
}

func (p *RespParser) Scan() (interface{}, error) {
	return nil, nil
}
