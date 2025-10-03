package command

import (
	"strings"
	"sync"

	"github.com/gretro/codecrafters-redis/app/protocol"
)

type Command interface {
	Name() string
	Execute(args *protocol.Array) (protocol.RespType, error)
}

var regInitOnce sync.Once
var registry *CommandRegistry

type CommandRegistry struct {
	commands map[string]Command
}

func (r *CommandRegistry) RegisterCommand(name string, command Command) {
	r.commands[strings.ToLower(name)] = command
}

func (r *CommandRegistry) Get(name string) Command {
	return r.commands[strings.ToLower(name)]
}

func (r *CommandRegistry) ResolveCommand(args *protocol.Array) (Command, error) {
	if len(args.Values) == 0 {
		return nil, &protocol.SimpleError{
			ErrorCode: "SYNTAX",
			Message:   "invalid command",
		}
	}

	commandName := args.Values[0].(*protocol.SimpleString).Value
	command := r.Get(commandName)

	if command == nil {
		return nil, &protocol.SimpleError{
			ErrorCode: "ERR",
			Message:   "unknown command",
		}
	}

	return command, nil
}

func GetRegistry() *CommandRegistry {
	regInitOnce.Do(func() {
		registry = &CommandRegistry{
			commands: make(map[string]Command),
		}

		registry.RegisterCommand("echo", &EchoCommand{})
		registry.RegisterCommand("ping", &PingCommand{})
	})

	return registry
}
