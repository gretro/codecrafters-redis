package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeSimpleError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{name: "error message", message: "error", expected: "-error\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.SimpleError{Message: test.message}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}
