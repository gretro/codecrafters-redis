package protocol

import (
	"bufio"
	"fmt"
	"io"
)

type RespParser struct {
	scanner *bufio.Scanner
}

func NewRespParser(r io.Reader) *RespParser {
	return &RespParser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *RespParser) Scan() (RespType, error) {
	ok := p.scanner.Scan()
	if !ok {
		return nil, io.EOF
	}

	firstLine := p.scanner.Text()
	if firstLine == "" {
		return nil, nil
	}

	identifier := string(firstLine[0])
	parser, ok := typeRegistry[identifier]
	if !ok {
		return nil, fmt.Errorf("%w: unknown first byte", ErrParse)
	}

	return parser(firstLine, p.scanner)
}
