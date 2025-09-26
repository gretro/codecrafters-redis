package protocol

import (
	"bufio"
	"strings"
	"testing"

	"github.com/gretro/codecrafters-redis/app/protocol"
	"github.com/stretchr/testify/require"
)

func Test_EncodeArray(t *testing.T) {
	tests := []struct {
		name     string
		value    []protocol.RespType
		expected string
	}{
		{name: "array of simple strings", value: []protocol.RespType{&protocol.SimpleString{Value: "1"}, &protocol.SimpleString{Value: "2"}, &protocol.SimpleString{Value: "3"}}, expected: "*3\r\n+1\r\n+2\r\n+3\r\n"},
		{name: "array of integers", value: []protocol.RespType{&protocol.Integer{Value: 1}, &protocol.Integer{Value: 2}, &protocol.Integer{Value: 3}}, expected: "*3\r\n:1\r\n:2\r\n:3\r\n"},
		{name: "array of booleans", value: []protocol.RespType{&protocol.Boolean{Value: true}, &protocol.Boolean{Value: false}, &protocol.Boolean{Value: true}}, expected: "*3\r\n#t\r\n#f\r\n#t\r\n"},
		{name: "array of bulk strings", value: []protocol.RespType{&protocol.BulkString{Value: "1"}, &protocol.BulkString{Value: "2"}, &protocol.BulkString{Value: "3"}}, expected: "*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n"},
		{name: "array of multi-line bulk strings", value: []protocol.RespType{&protocol.BulkString{Value: "1\r\n2\r\n3"}, &protocol.BulkString{Value: "4\r\n5\r\n6"}, &protocol.BulkString{Value: "7\r\n8\r\n9"}}, expected: "*3\r\n$7\r\n1\r\n2\r\n3\r\n$7\r\n4\r\n5\r\n6\r\n$7\r\n7\r\n8\r\n9\r\n"},
		{name: "array of mixed types", value: []protocol.RespType{&protocol.SimpleString{Value: "1"}, &protocol.Integer{Value: 2}, &protocol.Boolean{Value: true}}, expected: "*3\r\n+1\r\n:2\r\n#t\r\n"},
		{name: "empty array", value: []protocol.RespType{}, expected: "*0\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := &protocol.Array{Values: test.value}
			encoded := val.Encode()
			require.Equal(t, test.expected, string(encoded), "expected %s, got %s", test.expected, string(encoded))
		})
	}
}

func Test_ParseArray(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    []protocol.RespType
		expectedErr error
	}{
		{name: "array of simple strings", value: "*3\r\n+1\r\n+2\r\n+3\r\n", expected: []protocol.RespType{&protocol.SimpleString{Value: "1"}, &protocol.SimpleString{Value: "2"}, &protocol.SimpleString{Value: "3"}}, expectedErr: nil},
		{name: "array of integers", value: "*3\r\n:1\r\n:2\r\n:3\r\n", expected: []protocol.RespType{&protocol.Integer{Value: 1}, &protocol.Integer{Value: 2}, &protocol.Integer{Value: 3}}, expectedErr: nil},
		{name: "array of booleans", value: "*3\r\n#t\r\n#f\r\n#t\r\n", expected: []protocol.RespType{&protocol.Boolean{Value: true}, &protocol.Boolean{Value: false}, &protocol.Boolean{Value: true}}, expectedErr: nil},
		{name: "array of bulk strings", value: "*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n", expected: []protocol.RespType{&protocol.BulkString{Value: "1"}, &protocol.BulkString{Value: "2"}, &protocol.BulkString{Value: "3"}}, expectedErr: nil},
		{name: "array of multi-line bulk strings", value: "*3\r\n$7\r\n1\r\n2\r\n3\r\n$7\r\n4\r\n5\r\n6\r\n$7\r\n7\r\n8\r\n9\r\n", expected: []protocol.RespType{&protocol.BulkString{Value: "1\r\n2\r\n3"}, &protocol.BulkString{Value: "4\r\n5\r\n6"}, &protocol.BulkString{Value: "7\r\n8\r\n9"}}, expectedErr: nil},
		{name: "array of mixed types", value: "*3\r\n+1\r\n:2\r\n#t\r\n", expected: []protocol.RespType{&protocol.SimpleString{Value: "1"}, &protocol.Integer{Value: 2}, &protocol.Boolean{Value: true}}, expectedErr: nil},
		{name: "empty array", value: "*0\r\n", expected: []protocol.RespType{}, expectedErr: nil},
		{name: "invalid prefix", value: "+3\r\n+1\r\n+2\r\n+3\r\n", expected: nil, expectedErr: protocol.ErrParse},
		{name: "empty line in the middle", value: "*3\r\n+1\r\n\r\n+3\r\n", expected: nil, expectedErr: protocol.ErrParse},
		{name: "invalid length", value: "*invalid\r\n+1\r\n+2\r\n+3\r\n", expected: nil, expectedErr: protocol.ErrParse},
		{name: "negative length", value: "$-1\r\n+1\r\n+2\r\n+3\r\n", expected: nil, expectedErr: protocol.ErrParse},
		{name: "incorrect length", value: "*4\r\n+1\r\n+2\r\n+3\r\n", expected: nil, expectedErr: protocol.ErrParse},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(test.value))
			scanner.Scan()

			firstLine := scanner.Text()
			val, err := protocol.ParseArray(firstLine, scanner)

			if test.expectedErr == nil {
				require.NoError(t, err, "expected no error, got %v", err)
				require.IsType(t, &protocol.Array{}, val, "expected %T, got %T", protocol.Array{}, val)
				require.Equal(t, test.expected, val.(*protocol.Array).Values, "expected %v, got %v", test.expected, val.(*protocol.Array).Values)
			} else {
				require.ErrorIs(t, err, test.expectedErr, "expected error, got %v", err)
			}
		})
	}
}
