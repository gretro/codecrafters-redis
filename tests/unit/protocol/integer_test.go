package protocol

import (
	"fmt"
	"math"
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{name: "-1", value: -1, expected: ":-1\r\n"},
		{name: "1", value: 1, expected: ":1\r\n"},
		{name: "0", value: 0, expected: ":0\r\n"},
		{name: "min int64", value: math.MinInt64, expected: fmt.Sprintf(":%d\r\n", math.MinInt64)},
		{name: "max int64", value: math.MaxInt64, expected: fmt.Sprintf(":%d\r\n", math.MaxInt64)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.Integer{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseInteger(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    int64
		expectedErr error
	}{
		{name: "1", value: ":1", expected: 1, expectedErr: nil},
		{name: "0", value: ":0", expected: 0, expectedErr: nil},
		{name: "-1", value: ":-1", expected: -1, expectedErr: nil},
		{name: "-1 with EOL", value: ":-1\r\n", expected: -1, expectedErr: nil},
		{name: "+1", value: ":+1", expected: 1, expectedErr: nil},
		{name: "+1 with EOL", value: ":+1\r\n", expected: 1, expectedErr: nil},
		{name: "min int", value: fmt.Sprintf(":%d", math.MinInt64), expected: math.MinInt64, expectedErr: nil},
		{name: "max int", value: fmt.Sprintf(":%d", math.MaxInt64), expected: math.MaxInt64, expectedErr: nil},
		{name: "float", value: ":1.5", expected: 0, expectedErr: protocol.ErrParse},
		{name: "invalid prefix", value: "+1\r\n", expected: 0, expectedErr: protocol.ErrParse},
		{name: "invalid value", value: ":invalid\r\n", expected: 0, expectedErr: protocol.ErrParse},
		{name: "overflow", value: fmt.Sprintf(":%d0", math.MaxInt64), expected: 0, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseInteger([]byte(test.value))

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.Integer{}, val, "expected %T, got %T", protocol.Integer{}, val)
				require.Equal(t, test.expected, val.(*protocol.Integer).Value, "expected %d, got %d", test.expected, val.(*protocol.Integer).Value)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
