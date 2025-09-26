package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

type Integer struct {
	Value int64
}

const integerPrefix = ":"

func ParseInteger(data []byte) (RespType, error) {
	str := string(data)
	if !strings.HasPrefix(str, integerPrefix) {
		return nil, fmt.Errorf("%w: invalid integer prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(str, integerPrefix)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid integer value", ErrParse)
	}

	return &Integer{Value: value}, nil
}

func (i *Integer) Encode() []byte {
	return []byte(fmt.Sprintf("%s%d%s", integerPrefix, i.Value, EOL))
}
