package rpc

import (
	"net"
	"time"

	"github.com/dtynn/influxdbx/raftdemo/rpc/proto"
	"google.golang.org/grpc"
)

func NewClusterClient(target string, dialer func(string, time.Duration) (net.Conn, error)) (proto.ClusterClient, error) {
	cc, err := grpc.Dial(target, grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewClusterClient(cc), nil
}
