package run

import (
	"fmt"
	"net"
	"os"

	"github.com/dtynn/influxdbx/cluster"
	"github.com/dtynn/influxdbx/service/raft"
	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/logger"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxdb/tcp"
	"go.uber.org/zap"
)

// Service represents a service attached to the server.
type Service interface {
	WithLogger(log *zap.Logger)
	Open() error
	Close() error
}

type Server struct {
	MetaClient *meta.Client
	Cluster    *cluster.Cluster
	Raft       *raft.Raft

	logger *zap.Logger
	net.Listener

	tcpAddr           string
	reportingDisabled bool

	services []Service

	cfg *Config
}

func NewServer(cfg *Config) (*Server, error) {
	err := os.MkdirAll(cfg.Meta.Dir, 0777)
	if err != nil {
		return nil, err
	}

	_, err = influxdb.LoadNode(cfg.Meta.Dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	bind := cfg.BindAddress
	s := &Server{
		MetaClient: meta.NewClient(cfg.Meta),
		Raft:       raft.NewRaft(cfg.Meta.Dir, cfg.Raft),

		logger: logger.New(os.Stderr),

		tcpAddr:           bind,
		reportingDisabled: cfg.ReportingDisabled,

		cfg: cfg,
	}

	s.Cluster = cluster.NewCluster(cfg.Cluster)

	return s, nil
}

func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.tcpAddr)
	if err != nil {
		return err
	}

	s.logger.Info("listen on " + s.tcpAddr)

	s.Listener = ln
	mux := tcp.NewMux()
	go mux.Serve(ln)

	s.Raft.MetaClient = s.MetaClient
	s.Raft.Listener = mux.Listen(raft.MuxHeader)

	s.Cluster.NodeManager = s.Raft
	s.Cluster.Listener = mux.Listen(cluster.MuxHeader)

	s.appendRaftService()
	s.appendClusterService()

	for _, service := range s.services {
		if err := service.Open(); err != nil {
			return fmt.Errorf("open service error: %s", err)
		}
	}

	return nil
}

func (s *Server) Close() error {
	if s.Listener != nil {
		s.Listener.Close()
	}

	for _, service := range s.services {
		service.Close()
	}

	if s.MetaClient != nil {
		s.MetaClient.Close()
	}

	return nil
}

func (s *Server) appendRaftService() {
	if s.cfg.Raft.LoggingEnabled {
		s.Raft.WithLogger(s.logger)
	}

	s.services = append(s.services, s.Raft)
}

func (s *Server) appendClusterService() {
	s.services = append(s.services, s.Cluster)
}
