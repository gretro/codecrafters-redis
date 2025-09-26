package protocol

import (
	"bufio"
	"fmt"
	"strings"
)

type SimpleString struct {
	Value string
}

const SIMPLE_STRING_PREFIX = "+"

func ParseSimpleString(line string, _ *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(line, SIMPLE_STRING_PREFIX) {
		return nil, fmt.Errorf("%w: invalid simple string prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(line, SIMPLE_STRING_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	if strings.Contains(trimmed, "\r") || strings.Contains(trimmed, "\n") {
		return nil, fmt.Errorf("%w: simple string contains invalid characters", ErrParse)
	}

	return &SimpleString{Value: trimmed}, nil
}

func (s *SimpleString) Encode() []byte {
	return fmt.Appendf(nil, "%s%s%s", SIMPLE_STRING_PREFIX, s.Value, EOL)
}
