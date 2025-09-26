package protocol

import "errors"

var ErrParse = errors.New("failed to parse value")
var ErrUnknownTypeByte = errors.New("unknown type byte")

var EOL = "\r\n"
