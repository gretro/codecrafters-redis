package protocol

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeBulkString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "simple string", value: "Hello, world!", expected: "$13\r\nHello, world!\r\n"},
		{name: "empty string", value: "", expected: "$0\r\n\r\n"},
		{name: "multi-line string", value: "Hello, world!\r\nHello, world...", expected: "$30\r\nHello, world!\r\nHello, world...\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.BulkString{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseBulkString(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectedErr error
	}{
		{name: "simple string", value: "$13\r\nHello, world!\r\n", expected: "Hello, world!", expectedErr: nil},
		{name: "empty string", value: "$0\r\n\r\n", expected: "", expectedErr: nil},
		{name: "multi-line string", value: "$30\r\nHello, world!\r\nHello, world...\r\n", expected: "Hello, world!\r\nHello, world...", expectedErr: nil},
		{name: "invalid prefix", value: "+13\r\nHello, world!\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "invalid length", value: "$invalid\r\nHello, world!\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "negative length", value: "$-1\r\nHello, world!\r\n", expected: "", expectedErr: protocol.ErrParse},
		{name: "overflow", value: fmt.Sprintf("$%d\r\nHello, world!\r\n", protocol.MAX_BULK_STRING_LENGTH+1), expected: "", expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(test.value))
			scanner.Scan()

			firstLine := scanner.Text()
			val, err := protocol.ParseBulkString(firstLine, scanner)

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.BulkString{}, val, "expected %T, got %T", protocol.BulkString{}, val)
				require.Equal(t, test.expected, val.(*protocol.BulkString).Value, "expected %s, got %s", test.expected, val.(*protocol.BulkString).Value)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
