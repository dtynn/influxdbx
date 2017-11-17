package cluster

import (
	"context"
	"net"

	"github.com/dtynn/influxdbx/service/cluster/proto"
	"google.golang.org/grpc"
)

const (
	MuxHeader byte = 11
)

func NewCluster(impl proto.ClusterServer) *Cluster {
	s := grpc.NewServer()
	proto.RegisterClusterServer(s, impl)

	return &Cluster{
		s: s,
	}
}

type Cluster struct {
	net.Listener

	s *grpc.Server
}

func (c *Cluster) Open(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- c.s.Serve(c.Listener)
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		c.s.GracefulStop()

	case err := <-errCh:
		return err
	}

	return <-errCh
}
