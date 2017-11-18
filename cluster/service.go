package cluster

import (
	"net"

	"github.com/dtynn/influxdbx/cluster/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	MuxHeader byte = 11
)

func NewCluster(cfg *Config) *Cluster {
	c := &Cluster{
		s:      grpc.NewServer(),
		cfg:    cfg,
		logger: zap.NewNop(),
	}

	proto.RegisterClusterServer(c.s, c)

	return c
}

type Cluster struct {
	net.Listener

	s *grpc.Server

	NodeManager interface {
		AddNode(id, adress string) error
		RemoveNode(id string) error
		Leader() (string, bool)
		IsLeader() bool
		IsErrNotLeader(err error) bool
	}

	logger *zap.Logger

	cfg *Config
}

func (c *Cluster) Open() error {
	go c.s.Serve(c.Listener)
	return nil
}

func (c *Cluster) Close() error {
	c.s.GracefulStop()
	return nil
}

func (c *Cluster) WithLogger(l *zap.Logger) {
	c.logger = l.With(zap.String("service", "cluster"))
}
