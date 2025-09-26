package protocol

import (
	"bufio"
	"fmt"
	"maps"
	"sync"
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

func ResolvePrimitiveTypeParser(typeByte string) (TypeParser, error) {
	parser, ok := primitivesRegistry[typeByte]
	if !ok {
		return nil, fmt.Errorf("%w: unknown primitive type byte '%s'", ErrUnknownTypeByte, typeByte)
	}

	return parser, nil
}

var typeRegOnce sync.Once
var typeRegistry = map[string]TypeParser{}

func ResolveTypeParser(typeByte string) (TypeParser, error) {
	typeRegOnce.Do(func() {
		typeRegistry = maps.Clone(primitivesRegistry)
		typeRegistry[ARRAY_PREFIX] = ParseArray
	})

	parser, ok := typeRegistry[typeByte]
	if !ok {
		return nil, fmt.Errorf("%w: unknown type byte '%s'", ErrUnknownTypeByte, typeByte)
	}

	return parser, nil
}
