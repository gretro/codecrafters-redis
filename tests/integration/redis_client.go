package integration

import (
	"context"
	"net"
)

func NewTCPClient(ctx context.Context) (net.Conn, error) {
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
