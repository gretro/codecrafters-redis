package command

import "github.com/gretro/codecrafters-redis/app/protocol"

type PingCommand struct {
}

var _ Command = &PingCommand{}

func (c *PingCommand) Name() string {
	return "PING"
}

func (c *PingCommand) Execute(args *protocol.Array) (protocol.RespType, error) {
	return &protocol.SimpleString{Value: "PONG"}, nil
}
