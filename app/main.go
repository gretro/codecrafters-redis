package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/gretro/codecrafters-redis/app/command"
	"github.com/gretro/codecrafters-redis/app/logging"
	"github.com/gretro/codecrafters-redis/app/protocol"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting Redis server", "port", "6379")

	lc := net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", "0.0.0.0:6379")
	if err != nil {
		logger.Error("Failed to bind to port 6379", logging.ErrorKey, err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Failed to accept connection", logging.ErrorKey, err)
			continue
		}
		logger.Info("Accepted connection", logging.RemoteAddrKey, conn.RemoteAddr())
		go handleConnection(logger, conn)
	}
}

func handleConnection(l *slog.Logger, conn net.Conn) {
	defer func() {
		l.Info("Closing connection", logging.RemoteAddrKey, conn.RemoteAddr())

		err := conn.Close()
		if err != nil {
			l.Error("Failed to close connection", logging.ErrorKey, err)
		}
	}()

	parser := protocol.NewRespParser(conn)
	var err error
	var rawMsg protocol.RespType

	for ; err != io.EOF; rawMsg, err = parser.Scan() {
		startTime := time.Now()

		if err != nil {
			l.Error("Failed to parse command", logging.ErrorKey, err)

			var respErr protocol.RespType

			if !errors.As(err, &respErr) {
				respErr = &protocol.BulkError{
					ErrorCode: "ERR",
					Message:   err.Error(),
				}
			}

			_, err = conn.Write(respErr.Encode())

			if err != nil {
				l.Error("Failed to write error to connection", logging.ErrorKey, err, logging.RemoteAddrKey, conn.RemoteAddr())
				return
			}

			continue
		}

		msg, respErr := parseMessage(l, rawMsg)
		if respErr != nil {
			_, err = conn.Write(respErr.Encode())

			if err != nil {
				l.Error("Failed to write to TCP connection", logging.ErrorKey, err, logging.RemoteAddrKey, conn.RemoteAddr())
				return
			}

			continue
		}

		result, err := handleCommand(l, msg)
		if err != nil {
			if !errors.As(err, &respErr) {
				respErr = &protocol.BulkError{
					ErrorCode: "ERR",
					Message:   err.Error(),
				}
			}

			_, err = conn.Write(respErr.Encode())
			if err != nil {
				l.Error("Failed to write error to connection", logging.ErrorKey, err, logging.RemoteAddrKey, conn.RemoteAddr())
				return
			}

			continue
		}

		_, err = conn.Write(result.Result.Encode())
		if err != nil {
			l.Error("Failed to write response to connection", logging.ErrorKey, err, logging.RemoteAddrKey, conn.RemoteAddr())
			return
		}

		l.Info(fmt.Sprintf("Executed %s", result.CommandName), "duration_ms", time.Since(startTime).Milliseconds())
	}
}

func parseMessage(l *slog.Logger, msg protocol.RespType) (*protocol.Array, protocol.RespType) {
	msgArray, ok := msg.(*protocol.Array)
	if !ok {
		l.Error("Invalid message type. Expected *protocol.Array", "msgType", fmt.Sprintf("%T", msg))
		respErr := &protocol.BulkError{
			ErrorCode: "SYNTAX",
			Message:   "invalid message format",
		}

		return nil, respErr
	}

	return msgArray, nil
}

type CommandExecutionResult struct {
	CommandName string
	Result      protocol.RespType
}

func handleCommand(l *slog.Logger, msg *protocol.Array) (CommandExecutionResult, error) {
	command, err := command.GetRegistry().ResolveCommand(msg)
	if err != nil {
		l.Error("Failed to resolve command", logging.MsgKey, msg.String(), logging.ErrorKey, err)
		return CommandExecutionResult{
			CommandName: "UNKNOWN",
		}, err
	}

	resp, err := command.Execute(msg)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to execute %s", command.Name()), logging.MsgKey, msg.String(), logging.ErrorKey, err)
		return CommandExecutionResult{CommandName: command.Name()}, err
	}

	return CommandExecutionResult{
		CommandName: command.Name(),
		Result:      resp,
	}, nil
}
