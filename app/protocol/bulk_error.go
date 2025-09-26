package protocol

import (
	"fmt"
	"strings"
)

type BulkError struct {
	ErrorCode string
	Message   string
}

const BULK_ERROR_PREFIX = "!"

func (b *BulkError) Encode() []byte {
	msg := fmt.Sprintf("%s %s", strings.ToUpper(b.ErrorCode), b.Message)
	return fmt.Appendf(nil, "%s%d%s%s%s", BULK_ERROR_PREFIX, len(msg), EOL, msg, EOL)
}
