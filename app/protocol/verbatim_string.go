package protocol

import (
	"fmt"
)

type VerbatimString struct {
	Encoding string
	Value    string
}

const VERBATIM_STRING_PREFIX = "="

func (b *VerbatimString) Encode() []byte {
	enc := b.Encoding
	if len(enc) != 3 {
		panic("invalid VerbatimString: encoding must be 3 characters")
	}

	return fmt.Appendf(nil, "%s%d\r\n%s:%s\r\n", VERBATIM_STRING_PREFIX, len(b.Value), b.Encoding, b.Value)
}

func (b *VerbatimString) String() string {
	return fmt.Sprintf("(%s) `%s`", b.Encoding, b.Value)
}
