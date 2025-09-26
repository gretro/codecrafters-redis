package protocol

type RespType interface {
	Encode() []byte
}

type TypeParser func(data []byte) (RespType, error)

var typeRegistry = map[string]TypeParser{
	booleanPrefix:      ParseBoolean,
	integerPrefix:      ParseInteger,
	simpleStringPrefix: ParseSimpleString,
}
