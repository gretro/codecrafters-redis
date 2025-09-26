package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeBoolean(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{name: "true", value: true, expected: "#t\r\n"},
		{name: "false", value: false, expected: "#f\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.Boolean{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseBoolean(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    bool
		expectedErr error
	}{
		{name: "true", value: "#t", expected: true, expectedErr: nil},
		{name: "false", value: "#f", expected: false, expectedErr: nil},
		{name: "true with EOL", value: "#t\r\n", expected: true, expectedErr: nil},
		{name: "false with EOL", value: "#f\r\n", expected: false, expectedErr: nil},
		{name: "invalid prefix", value: "+t\r\n", expected: false, expectedErr: protocol.ErrParse},
		{name: "invalid value", value: "#invalid\r\n", expected: false, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseBoolean([]byte(test.value))

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.Boolean{}, val, "expected %T, got %T", protocol.Boolean{}, val)
				require.Equal(t, test.expected, val.(*protocol.Boolean).Value, "expected %v, got %v", test.expected, val.(*protocol.Boolean).Value)
			} else {
				require.Error(t, err, "expected error, got %v", err)
			}
		})
	}
}
