package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Integer struct {
	Value int64
}

const INTEGER_PREFIX = ":"

func ParseInteger(line string, _ *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(line, INTEGER_PREFIX) {
		return nil, fmt.Errorf("%w: invalid integer prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(line, INTEGER_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid integer value", ErrParse)
	}

	return &Integer{Value: value}, nil
}

func (i *Integer) Encode() []byte {
	return fmt.Appendf(nil, "%s%d%s", INTEGER_PREFIX, i.Value, EOL)
}
