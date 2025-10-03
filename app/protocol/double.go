package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Double struct {
	Value float64
}

const DOUBLE_PREFIX = ","

func ParseDouble(firstLine string, scanner *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(firstLine, DOUBLE_PREFIX) {
		return nil, fmt.Errorf("%w: invalid double prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(firstLine, DOUBLE_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid double value", ErrParse)
	}

	return &Double{Value: value}, nil
}

func (d *Double) Encode() []byte {
	return fmt.Appendf(nil, "%s%G%s", DOUBLE_PREFIX, d.Value, EOL)
}

func (d *Double) String() string {
	return fmt.Sprintf("Float64(%G)", d.Value)
}
