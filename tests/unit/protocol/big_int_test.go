package protocol

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeBigInt(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		expected string
	}{
		{name: "1", value: big.NewInt(1), expected: "(1\r\n"},
		{name: "0", value: big.NewInt(0), expected: "(0\r\n"},
		{name: "-1", value: big.NewInt(-1), expected: "(-1\r\n"},
		{name: "bigger than max int64", value: func() *big.Int {
			bigInt, ok := big.NewInt(0).SetString(fmt.Sprintf("%d00", math.MaxInt64), 10)
			if !ok {
				t.Fatal("failed to set string")
			}
			return bigInt
		}(), expected: fmt.Sprintf("(%d00\r\n", math.MaxInt64)},
		{name: "lower than min int64", value: func() *big.Int {
			bigInt, ok := big.NewInt(0).SetString(fmt.Sprintf("%d99", math.MinInt64), 10)
			if !ok {
				t.Fatal("failed to set string")
			}
			return bigInt
		}(), expected: fmt.Sprintf("(%d99\r\n", math.MinInt64)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.BigInt{Value: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseBigInt(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    *big.Int
		expectedErr error
	}{
		{name: "1", value: "(1\r\n", expected: big.NewInt(1), expectedErr: nil},
		{name: "0", value: "(0\r\n", expected: big.NewInt(0), expectedErr: nil},
		{name: "-1", value: "(-1\r\n", expected: big.NewInt(-1), expectedErr: nil},
		{name: "max int64", value: fmt.Sprintf("(%d\r\n", math.MaxInt64), expected: big.NewInt(math.MaxInt64), expectedErr: nil},
		{name: "min int64", value: fmt.Sprintf("(%d\r\n", math.MinInt64), expected: big.NewInt(math.MinInt64), expectedErr: nil},
		{name: "bigger than max int64", value: fmt.Sprintf("(%d00\r\n", math.MaxInt64), expected: func() *big.Int {
			bigInt, ok := big.NewInt(0).SetString(fmt.Sprintf("%d00", math.MaxInt64), 10)
			if !ok {
				t.Fatal("failed to set string")
			}
			return bigInt
		}(), expectedErr: nil},
		{name: "lower than min int64", value: fmt.Sprintf("(%d99\r\n", math.MinInt64), expected: func() *big.Int {
			bigInt, ok := big.NewInt(0).SetString(fmt.Sprintf("%d99", math.MinInt64), 10)
			if !ok {
				t.Fatal("failed to set string")
			}

			return bigInt
		}(), expectedErr: nil},
		{name: "invalid prefix", value: "+1\r\n", expected: nil, expectedErr: protocol.ErrParse},
		{name: "invalid value", value: "(invalid\r\n", expected: nil, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := protocol.ParseBigInt(test.value, nil)
			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.BigInt{}, val, "expected %T, got %T", protocol.BigInt{}, val)
				require.Equal(t, test.expected, val.(*protocol.BigInt).Value, "expected %d, got %d", test.expected, val.(*protocol.BigInt).Value)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
