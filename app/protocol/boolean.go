package protocol

import (
	"bufio"
	"fmt"
	"strings"
)

const BOOLEAN_PREFIX = "#"

type Boolean struct {
	Value bool
}

func ParseBoolean(line string, _ *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(line, BOOLEAN_PREFIX) {
		return nil, fmt.Errorf("%w: invalid boolean prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(line, BOOLEAN_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	isValid := false
	value := false
	if strings.ToLower(trimmed) == "t" {
		value = true
		isValid = true
	} else if strings.ToLower(trimmed) == "f" {
		value = false
		isValid = true
	}

	if !isValid {
		return nil, fmt.Errorf("%w: invalid boolean value", ErrParse)
	}

	return &Boolean{Value: value}, nil
}

func (b *Boolean) Encode() []byte {
	strVal := ""
	if b.Value {
		strVal = "t"
	} else {
		strVal = "f"
	}

	return fmt.Appendf(nil, "%s%s%s", BOOLEAN_PREFIX, strVal, EOL)
}

func (b *Boolean) String() string {
	return fmt.Sprintf("Boolean(%v)", b.Value)
}
