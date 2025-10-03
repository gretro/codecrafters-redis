package protocol

import (
	"bufio"
	"fmt"
	"math/big"
	"strings"
)

type BigInt struct {
	Value *big.Int
}

const BIG_NUMBER_PREFIX = "("

func ParseBigInt(firstLine string, scanner *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(firstLine, BIG_NUMBER_PREFIX) {
		return nil, fmt.Errorf("%w: invalid big number prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(firstLine, BIG_NUMBER_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	value, ok := big.NewInt(0).SetString(trimmed, 10)
	if !ok {
		return nil, fmt.Errorf("%w: invalid big number value", ErrParse)
	}

	return &BigInt{Value: value}, nil
}

func (b *BigInt) Encode() []byte {
	return fmt.Appendf(nil, "%s%s%s", BIG_NUMBER_PREFIX, b.Value.Text(10), EOL)
}

func (b *BigInt) String() string {
	return fmt.Sprintf("BigInt(%s)", b.Value.Text(10))
}
