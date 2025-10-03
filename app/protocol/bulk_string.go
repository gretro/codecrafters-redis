package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type BulkString struct {
	Value string
}

const BULK_STRING_PREFIX = "$"
const MAX_BULK_STRING_LENGTH = 512 * 1024 * 1024 // 512MB

func ParseBulkString(firstLine string, scanner *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(firstLine, BULK_STRING_PREFIX) {
		return nil, fmt.Errorf("%w: invalid bulk string prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(firstLine, BULK_STRING_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	expectedLength, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid bulk string length", ErrParse)
	}

	if expectedLength > MAX_BULK_STRING_LENGTH {
		return nil, fmt.Errorf("%w: bulk string length exceeds maximum length", ErrParse)
	}

	actualLength := uint64(0)
	strBuilder := strings.Builder{}

	for actualLength < expectedLength {
		ok := scanner.Scan()

		line := scanner.Text()
		strBuilder.WriteString(line)
		actualLength += uint64(len(line))

		// If the actual length is less than the expected length, that means we have a multi-line string.
		if actualLength < expectedLength {
			if !ok {
				return nil, fmt.Errorf("%w: premature end of stream. Expected %d bytes, got %d bytes", ErrParse, expectedLength, actualLength)
			}

			strBuilder.WriteString(EOL)
			actualLength += uint64(len(EOL))
		}
	}

	return &BulkString{Value: strBuilder.String()}, nil
}

func (b *BulkString) Encode() []byte {
	return fmt.Appendf(nil, "%s%d\r\n%s\r\n", BULK_STRING_PREFIX, len(b.Value), b.Value)
}

func (b *BulkString) String() string {
	return fmt.Sprintf("`%s`", b.Value)
}
