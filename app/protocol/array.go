package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Array struct {
	Values []RespType
}

const ARRAY_PREFIX = "*"

func ParseArray(firstLine string, scanner *bufio.Scanner) (RespType, error) {
	if !strings.HasPrefix(firstLine, ARRAY_PREFIX) {
		return nil, fmt.Errorf("%w: invalid array prefix", ErrParse)
	}

	trimmed := strings.TrimPrefix(firstLine, ARRAY_PREFIX)
	trimmed = strings.TrimSuffix(trimmed, EOL)

	expectedLength, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid array length", ErrParse)
	}

	actualLength := uint64(0)
	values := make([]RespType, expectedLength)
	errs := make([]error, 0)

	for actualLength < expectedLength {
		scanOK := scanner.Scan()

		line := scanner.Text()
		if len(line) == 0 {
			errs = append(errs, fmt.Errorf("%w: received invalid empty value in array at line %d", ErrParse, actualLength))
			actualLength++
			continue
		}

		typeByte := line[0]
		parser, err := ResolvePrimitiveTypeParser(string(typeByte))
		if err != nil {
			return nil, fmt.Errorf("%w: error parsing value at line %d: %w", ErrParse, actualLength, err)
		}

		value, err := parser(line, scanner)
		if err != nil {
			errs = append(errs, fmt.Errorf("%w: failed to parse value in array at line %d", ErrParse, actualLength))
			actualLength++
		}

		values[actualLength] = value
		actualLength++

		if actualLength < expectedLength && !scanOK {
			return nil, fmt.Errorf("%w: premature end of stream. Expected %d lines, got %d lines", ErrParse, expectedLength, actualLength)
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse values in array: %w", errors.Join(errs...))
	}

	return &Array{
		Values: values,
	}, nil
}

func (a *Array) Encode() []byte {
	lengthLine := fmt.Sprintf("%s%d\r\n", ARRAY_PREFIX, len(a.Values))

	encodedValues := strings.Builder{}
	for _, value := range a.Values {
		encodedValues.Write(value.Encode())
	}

	return fmt.Appendf(nil, "%s%s", lengthLine, encodedValues.String())
}

func (a *Array) String() string {
	vals := make([]string, len(a.Values))
	for i, value := range a.Values {
		vals[i] = fmt.Sprintf("%s", value.String())
	}

	joined := strings.Join(vals, ", ")

	return fmt.Sprintf("[%s]", joined)
}
