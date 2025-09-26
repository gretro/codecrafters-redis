package protocol

import "fmt"

type SimpleError struct {
	Message string
}

const simpleErrorPrefix = "-"

func (s *SimpleError) Encode() []byte {
	return []byte(fmt.Sprintf("%s%s%s", simpleErrorPrefix, s.Message, EOL))
}
