package protocol

import (
	"bufio"
	"fmt"
	"strings"
)

type Null struct {
}

const NULL_PREFIX = "_"

func ParseNull(firstLine string, scanner *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(firstLine, NULL_PREFIX) {
		return nil, fmt.Errorf("%w: invalid null prefix", ErrParse)
	}

	return &Null{}, nil
}

func (n *Null) Encode() []byte {
	return fmt.Appendf(nil, "%s%s", NULL_PREFIX, EOL)
}
