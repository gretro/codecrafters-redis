package command

import "github.com/gretro/codecrafters-redis/app/protocol"

type EchoCommand struct {
}

var _ Command = &EchoCommand{}

func (c *EchoCommand) Name() string {
	return "ECHO"
}

func (c *EchoCommand) extractArgs(args *protocol.Array) (string, error) {
	if len(args.Values) != 2 {
		return "", &protocol.SimpleError{
			ErrorCode: "SYNTAX",
			Message:   "invalid number of arguments",
		}
	}

	msg := args.Values[1]

	if simpleStr, ok := msg.(*protocol.SimpleString); ok {
		return simpleStr.Value, nil
	}

	if bulkStr, ok := msg.(*protocol.BulkString); ok {
		return bulkStr.Value, nil
	}

	return "", &protocol.SimpleError{
		ErrorCode: "SYNTAX",
		Message:   "invalid argument type",
	}
}

func (c *EchoCommand) Execute(args *protocol.Array) (protocol.RespType, error) {
	msg, err := c.extractArgs(args)
	if err != nil {
		return nil, err
	}

	return &protocol.BulkString{Value: msg}, nil
}
