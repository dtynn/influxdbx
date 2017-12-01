package tcputil

import (
	"fmt"
	"net"
	"time"
)

// Dial connects to a remote mux listener with a given header byte.
func Dial(network, address string, timeout time.Duration, header byte) (net.Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Write([]byte{header}); err != nil {
		return nil, fmt.Errorf("write mux header: %s", err)
	}

	return conn, nil
}

