package protocol

import (
	"fmt"
	"strings"
)

type SimpleError struct {
	ErrorCode string
	Message   string
}

const SIMPLE_ERROR_PREFIX = "-"

func (s *SimpleError) Encode() []byte {
	return []byte(fmt.Sprintf("%s%s %s%s", SIMPLE_ERROR_PREFIX, strings.ToUpper(s.ErrorCode), s.Message, EOL))
}
