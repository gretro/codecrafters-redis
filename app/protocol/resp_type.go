package protocol

import (
	"bufio"
	"maps"
)

type RespType interface {
	Encode() []byte
}

type TypeParser func(firstLine string, scanner *bufio.Scanner) (RespType, error)

var primitivesRegistry = map[string]TypeParser{
	BOOLEAN_PREFIX:       ParseBoolean,
	INTEGER_PREFIX:       ParseInteger,
	SIMPLE_STRING_PREFIX: ParseSimpleString,
	BULK_STRING_PREFIX:   ParseBulkString,
}

var typeRegistry = map[string]TypeParser{}

func init() {
	maps.Copy(typeRegistry, primitivesRegistry)
}
