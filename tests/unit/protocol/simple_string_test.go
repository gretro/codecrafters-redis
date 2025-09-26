package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "1", value: "1", expected: "+1\r\n"},
		{name: "0", value: "0", expected: "+0\r\n"},
		{name: "OK", value: "OK", expected: "+OK\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.SimpleString{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseSimpleString(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectedErr error
	}{
		{name: "1", value: "+1\r\n", expected: "1", expectedErr: nil},
		{name: "0", value: "+0\r\n", expected: "0", expectedErr: nil},
		{name: "OK", value: "+OK\r\n", expected: "OK", expectedErr: nil},
		{name: "1 with EOL", value: "+1\r\n", expected: "1", expectedErr: nil},
		{name: "0 with EOL", value: "+0\r\n", expected: "0", expectedErr: nil},
		{name: "OK with EOL", value: "+OK\r\n", expected: "OK", expectedErr: nil},
		{name: "empty value", value: "+\r\n", expected: "", expectedErr: nil},
		{name: "invalid prefix", value: "-1\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "invalid value with CR", value: "+\r\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "invalid value with LF", value: "+\n\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "invalid value with CRLF", value: "+\r\n\r\n", expected: "", expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseSimpleString([]byte(test.value))

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.SimpleString{}, val, "expected %T, got %T", protocol.SimpleString{}, val)
				require.Equal(t, test.expected, val.(*protocol.SimpleString).Value, "expected %s, got %s", test.expected, val.(*protocol.SimpleString).Value)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
