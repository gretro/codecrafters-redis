package integration

import (
	"net"
)

func NewTcpClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
