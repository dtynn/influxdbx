package service

import (
	"net"
	"time"

	"github.com/dtynn/influxdbx/tcp"
	"github.com/hashicorp/raft"
)

type demoStreamLayer struct {
	net.Listener
}

func (d *demoStreamLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	return tcp.Dial("tcp", string(address), timeout, RAFTHeader)
}
