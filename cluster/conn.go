package cluster

import (
	"net"
	"sync"
	"time"

	"github.com/dtynn/influxdbx/tcp"
	"google.golang.org/grpc"
)

var (
	conns   = map[string]*grpc.ClientConn{}
	connsMu sync.Mutex
)

func newConn(target string) (*grpc.ClientConn, error) {
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		return tcp.Dial("tcp", address, timeout, MuxHeader)
	}

	cc, err := grpc.Dial(target, grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func getConn(target string) (*grpc.ClientConn, error) {
	connsMu.Lock()
	defer connsMu.Unlock()

	cc, ok := conns[target]
	if ok {
		return cc, nil
	}

	cc, err := newConn(target)
	if err != nil {
		return nil, err
	}

	conns[target] = cc

	return cc, nil
}
