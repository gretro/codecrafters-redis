package protocol

import (
	"fmt"
	"strings"
)

type SimpleError struct {
	ErrorCode string
	Message   string
}

var _ error = &SimpleError{}

const SIMPLE_ERROR_PREFIX = "-"

func (s *SimpleError) Encode() []byte {
	return []byte(fmt.Sprintf("%s%s %s%s", SIMPLE_ERROR_PREFIX, strings.ToUpper(s.ErrorCode), s.Message, EOL))
}

func (s *SimpleError) Error() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(s.ErrorCode), s.Message)
}

func (s *SimpleError) String() string {
	return fmt.Sprintf("SimpleError(%s %s)", strings.ToUpper(s.ErrorCode), s.Message)
}
