package protocol

import (
	"fmt"
	"strings"
)

type SimpleString struct {
	Value string
}

const simpleStringPrefix = "+"

func ParseSimpleString(data []byte) (RespType, error) {
	str := string(data)

	if !strings.HasPrefix(str, simpleStringPrefix) {
		return nil, fmt.Errorf("%w: invalid simple string prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(str, simpleStringPrefix)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	if strings.Contains(trimmed, "\r") || strings.Contains(trimmed, "\n") {
		return nil, fmt.Errorf("%w: simple string contains invalid characters", ErrParse)
	}

	return &SimpleString{Value: trimmed}, nil
}

func (s *SimpleString) Encode() []byte {
	return []byte(fmt.Sprintf("%s%s%s", simpleStringPrefix, s.Value, EOL))
}
