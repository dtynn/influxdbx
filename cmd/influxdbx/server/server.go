package server

import (
	"fmt"
	"net"
	"os"

	"github.com/dtynn/influxdbx/service/meta"
	"github.com/dtynn/influxdbx/service/rpc"
	"github.com/influxdata/influxdb/logger"
	"github.com/influxdata/influxdb/tcp"
	"go.uber.org/zap"
)

// Service represents a service attached to the server.
type Service interface {
	Open() error
	Close() error
	WithLogger(log *zap.Logger)
}

// Server influxdbx server
type Server struct {
	Meta *meta.Meta
	RPC  *rpc.RPC

	logger *zap.Logger
	net.Listener

	bindAddress       string
	reportingDisabled bool

	services []Service

	cfg *Config
}

// NewServer return new server
func NewServer(cfg *Config) (*Server, error) {
	err := os.MkdirAll(cfg.Dir, 0777)
	if err != nil {
		return nil, err
	}

	bind := cfg.BindAddress
	s := &Server{
		RPC: rpc.NewRPC(cfg.RPC),

		logger: logger.New(os.Stderr),

		bindAddress:       bind,
		reportingDisabled: cfg.ReportingDisabled,

		cfg: cfg,
	}

	if cfg.Meta.Enabled {
		s.Meta = meta.NewMeta(cfg.Dir, cfg.Meta)
		s.RPC.MetaNodeManager = s.Meta
	}

	return s, nil
}

// Open open the server
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}

	s.logger.Info("listen on " + s.bindAddress)

	s.Listener = ln
	mux := tcp.NewMux()
	go mux.Serve(ln)

	s.RPC.Listener = mux.Listen(rpc.MuxHeader)

	if s.Meta != nil {
		s.Meta.Listener = mux.Listen(meta.MuxHeader)
	}

	s.appendRPCService()
	s.appendMetaService()

	for _, service := range s.services {
		if err := service.Open(); err != nil {
			return fmt.Errorf("open service error: %s", err)
		}
	}

	return nil
}

// Close close the server
func (s *Server) Close() error {
	if s.Listener != nil {
		s.Listener.Close()
	}

	for _, service := range s.services {
		service.Close()
	}

	return nil
}

func (s *Server) appendRPCService() {
	if s.cfg.RPC.Verbose {
		s.RPC.WithLogger(s.logger)
	}

	s.services = append(s.services, s.RPC)
}

func (s *Server) appendMetaService() {
	if !s.cfg.Meta.Enabled {
		return
	}

	if s.cfg.Meta.Verbose {
		s.Meta.WithLogger(s.logger)
	}

	s.services = append(s.services, s.Meta)
}
