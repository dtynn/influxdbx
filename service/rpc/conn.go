package rpc

import (
	"net"
	"sync"
	"time"

	"github.com/dtynn/influxdbx/util/tcputil"
	"google.golang.org/grpc"
)

var (
	defaultConnMgr = NewConnectionMgr()
)

// NewConn return new grpc connection
func NewConn(target string) (*grpc.ClientConn, error) {
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		return tcputil.Dial("tcp", address, timeout, MuxHeader)
	}

	cc, err := grpc.Dial(target, grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return cc, nil

}

// NewConnectionMgr return a new connection mgr
func NewConnectionMgr() *ConnectionMgr {
	return &ConnectionMgr{
		conns: map[string]*grpc.ClientConn{},
	}
}

// ConnectionMgr rpc connection manager
type ConnectionMgr struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

// Get get a connection
func (c *ConnectionMgr) Get(target string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cc, ok := c.conns[target]
	if ok {
		return cc, nil
	}

	cc, err := NewConn(target)
	if err != nil {
		return nil, err
	}

	c.conns[target] = cc

	return cc, nil
}
