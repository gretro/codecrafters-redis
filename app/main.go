package main

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"os"
	"strings"
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
		logger.Error("Failed to bind to port 6379", "error", err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Failed to accept connection", "error", err)
			continue
		}
		logger.Info("Accepted connection", "address", conn.RemoteAddr())
		go handleConnection(logger, conn)
	}
}

func handleConnection(l *slog.Logger, conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			l.Error("Failed to close connection", "error", err)
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		l.Info("Received command", "command", command)
		if strings.TrimSpace(command) == "PING" {
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				l.Error("Failed to write to connection", "error", err)
				return
			}
		} else {
			l.Warn("Unknown command", "command", command)
			// conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}

	l.Info("Closed connection", "address", conn.RemoteAddr())
}
