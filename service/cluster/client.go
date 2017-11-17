package cluster

import (
	"net"
	"time"

	"github.com/dtynn/influxdbx/service/cluster/proto"
	"github.com/dtynn/influxdbx/tcp"
	"google.golang.org/grpc"
)

func NewClient(target string) (proto.ClusterClient, error) {
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		return tcp.Dial("tcp", address, timeout, MuxHeader)
	}

	cc, err := grpc.Dial(target, grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewClusterClient(cc), nil
}
