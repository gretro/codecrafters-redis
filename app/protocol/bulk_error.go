package protocol

import (
	"fmt"
	"strings"
)

type BulkError struct {
	ErrorCode string
	Message   string
}

var _ error = &BulkError{}

const BULK_ERROR_PREFIX = "!"

func (b *BulkError) Encode() []byte {
	msg := fmt.Sprintf("%s %s", strings.ToUpper(b.ErrorCode), b.Message)
	return fmt.Appendf(nil, "%s%d%s%s%s", BULK_ERROR_PREFIX, len(msg), EOL, msg, EOL)
}

func (b *BulkError) Error() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(b.ErrorCode), b.Message)
}

func (b *BulkError) String() string {
	return fmt.Sprintf("BulkError(%s %s)", strings.ToUpper(b.ErrorCode), b.Message)
}
