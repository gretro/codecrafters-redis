package protocol

import (
	"fmt"
	"strings"
)

const booleanPrefix = "#"

type Boolean struct {
	Value bool
}

func ParseBoolean(data []byte) (RespType, error) {
	str := string(data)

	if !strings.HasPrefix(str, booleanPrefix) {
		return nil, fmt.Errorf("%w: invalid boolean prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(str, booleanPrefix)
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

	return []byte(fmt.Sprintf("%s%s%s", booleanPrefix, strVal, EOL))
}
