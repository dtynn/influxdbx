package rpc

import (
	"context"
	"log"
	"net"

	"github.com/dtynn/influxdbx/raftdemo/rpc/proto"
	"google.golang.org/grpc"
)

func NewCluster(ln net.Listener, srv proto.ClusterServer) *ClusterService {
	s := grpc.NewServer()
	proto.RegisterClusterServer(s, srv)

	return &ClusterService{
		ln: ln,
		s:  s,
	}
}

type ClusterService struct {
	ln net.Listener
	s  *grpc.Server
}

func (c *ClusterService) Run(ctx context.Context) error {
	log.Println("cluster service on", c.ln.Addr().String())

	errCh := make(chan error, 1)

	go func() {
		errCh <- c.s.Serve(c.ln)
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
