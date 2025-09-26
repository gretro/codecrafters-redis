package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeBulkError(t *testing.T) {
	tests := []struct {
		name      string
		errorCode string
		message   string
		expected  string
	}{
		{name: "error", errorCode: "SYNTAX", message: "error", expected: "!12\r\nSYNTAX error\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.BulkError{ErrorCode: test.errorCode, Message: test.message}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}
