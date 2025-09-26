package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeNull(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{name: "null", expected: "_\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.Null{}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseNull(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    *protocol.Null
		expectedErr error
	}{
		{name: "null", value: "_", expected: &protocol.Null{}, expectedErr: nil},
		{name: "null with EOL", value: "_\r\n", expected: &protocol.Null{}, expectedErr: nil},
		{name: "invalid prefix", value: "+_", expected: nil, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseNull(test.value, nil)

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.Null{}, val, "expected %T, got %T", protocol.Null{}, val)
				require.Equal(t, test.expected, val.(*protocol.Null), "expected %v, got %v", test.expected, val.(*protocol.Null))
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
