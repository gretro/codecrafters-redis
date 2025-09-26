package protocol

import "bufio"

type RespType interface {
	Encode() []byte
}

type TypeParser func(firstLine string, scanner *bufio.Scanner) (RespType, error)

var typeRegistry = map[string]TypeParser{
	BOOLEAN_PREFIX:       ParseBoolean,
	INTEGER_PREFIX:       ParseInteger,
	SIMPLE_STRING_PREFIX: ParseSimpleString,
}
