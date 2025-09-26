package main

import (
	"bufio"
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
	l, err := net.Listen("tcp", "0.0.0.0:6379")
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		if strings.TrimSpace(command) == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
