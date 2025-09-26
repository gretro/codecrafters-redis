package integration

import (
	"bufio"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	client, err := NewTcpClient()
	require.NoError(t, err, "failed to create TCP client")

	defer client.Close()

	_, err = client.Write([]byte("PING\r\n"))
	require.NoError(t, err, "failed to write to TCP client")

	scanner := bufio.NewScanner(client)
	scanner.Scan()
	require.NoError(t, scanner.Err(), "failed to read from TCP client")
	require.Equal(t, "+PONG", scanner.Text(), "expected PONG response")
}

func TestManyPings(t *testing.T) {
	client, err := NewTcpClient()
	require.NoError(t, err, "failed to create TCP client")

	defer client.Close()

	_, err = client.Write([]byte("PING\r\nPING\r\n"))
	require.NoError(t, err, "failed to write to TCP client")

	scanner := bufio.NewScanner(client)
	scanner.Scan()
	require.NoError(t, scanner.Err(), "failed to read from TCP client")
	require.Equal(t, "+PONG", scanner.Text(), "expected PONG response")

	scanner.Scan()
	require.NoError(t, scanner.Err(), "failed to read from TCP client")
	require.Equal(t, "+PONG", scanner.Text(), "expected PONG response")
}
