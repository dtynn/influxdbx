package raft

import (
	"net"
	"time"

	"github.com/dtynn/influxdbx/tcp"
	"github.com/hashicorp/raft"
)

type Stream struct {
	net.Listener
}

func (s *Stream) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	return tcp.Dial("tcp", string(address), timeout, MuxHeader)
}

func NewStream(ln net.Listener) *Stream {
	return &Stream{
		Listener: ln,
	}
}
