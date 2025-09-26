package protocol

import "fmt"

type SimpleError struct {
	Message string
}

const SIMPLE_ERROR_PREFIX = "-"

func (s *SimpleError) Encode() []byte {
	return []byte(fmt.Sprintf("%s%s%s", SIMPLE_ERROR_PREFIX, s.Message, EOL))
}
