package meta

import (
	"net"
	"time"

	"github.com/dtynn/influxdbx/util/tcputil"
	"github.com/hashicorp/raft"
)

var (
	_ raft.StreamLayer = (*stream)(nil)
)

func newStream(ln net.Listener) *stream {
	return &stream{
		Listener: ln,
	}
}

type stream struct {
	net.Listener
}

func (s *stream) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	return tcputil.Dial("tcp", string(address), timeout, MuxHeader)
}
