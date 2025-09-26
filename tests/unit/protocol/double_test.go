package protocol

import (
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeDouble(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{name: "1", value: 1.0, expected: ",1\r\n"},
		{name: "1000", value: 1000, expected: ",1000\r\n"},
		{name: "1_000_000", value: 1_000_000, expected: ",1E+06\r\n"},
		{name: "3.45E-06", value: 0.00000345, expected: ",3.45E-06\r\n"},
		{name: "-3.1416", value: -3.1416, expected: ",-3.1416\r\n"},
		{name: "-0.00003567", value: -0.00003567, expected: ",-3.567E-05\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.Double{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseDouble(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    float64
		expectedErr error
	}{
		{name: "1", value: ",1\r\n", expected: 1, expectedErr: nil},
		{name: "1000", value: ",1000\r\n", expected: 1000, expectedErr: nil},
		{name: "1E+06", value: ",1E+06\r\n", expected: 1000000, expectedErr: nil},
		{name: "3.45e-06", value: ",3.45e-06\r\n", expected: 0.00000345, expectedErr: nil},
		{name: "-3.1416", value: ",-3.1416\r\n", expected: -3.1416, expectedErr: nil},
		{name: "-3.567E-05", value: ",-3.567E-05\r\n", expected: -0.00003567, expectedErr: nil},
		{name: "invalid prefix", value: "+1\r\n", expected: 0, expectedErr: protocol.ErrParse},
		{name: "invalid value", value: ",invalid\r\n", expected: 0, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseDouble(test.value, nil)

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.Double{}, val, "expected %T, got %T", protocol.Double{}, val)
				require.Equal(t, test.expected, val.(*protocol.Double).Value, "expected %f, got %f", test.expected, val.(*protocol.Double).Value)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
