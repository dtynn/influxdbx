package rpc

import (
	"net"

	"github.com/dtynn/influxdbx/internal/pb"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// MuxHeader tcp mux header
	MuxHeader byte = 7
)

var (
	_ pb.MetaServer = (*RPC)(nil)
)

// NewRPC return rpc service
func NewRPC(cfg *Config) *RPC {
	return &RPC{
		logger: zap.NewNop(),
	}
}

// RPC rpc service
type RPC struct {
	s *grpc.Server

	net.Listener

	MetaNodeManager interface {
		AddMetaNode(id uint64, address string) error
		RemoveMetaNode(id uint64) error
		IsLeader() bool
		Leader() string
	}

	logger *zap.Logger
}

// Open open the rpc service
func (r *RPC) Open() error {
	r.s = grpc.NewServer()

	if r.MetaNodeManager != nil {
		pb.RegisterMetaServer(r.s, r)
	}

	go r.s.Serve(r.Listener)
	return nil
}

// Close close the rpc service
func (r *RPC) Close() error {
	if r.s != nil {
		r.s.GracefulStop()
	}

	return nil
}

// WithLogger setup logger
func (r *RPC) WithLogger(l *zap.Logger) {
	r.logger = l.With(zap.String("service", "rpc"))
}

// Join add meta node
func (r *RPC) Join(ctx context.Context, req *pb.MetaJoinReq) (*pb.MetaJoinResp, error) {
	err := r.MetaNodeManager.AddMetaNode(req.GetId(), req.GetAddress())
	resp := &pb.MetaJoinResp{
		Result: &pb.ClusterResult{},
	}

	resp.Result.Code, resp.Result.Msg = Err2Code(err)
	if canForward(resp.Result.Code) {
		leader := r.MetaNodeManager.Leader()
		resp.Result.Leader = &pb.ClusterLeader{
			Available: leader != "",
			Address:   leader,
		}
	}

	return resp, nil
}

// Remove remove meta node
func (r *RPC) Remove(ctx context.Context, req *pb.MetaRemoveReq) (*pb.MetaRemoveResp, error) {
	err := r.MetaNodeManager.RemoveMetaNode(req.GetId())
	resp := &pb.MetaRemoveResp{
		Result: &pb.ClusterResult{},
	}

	resp.Result.Code, resp.Result.Msg = Err2Code(err)
	if canForward(resp.Result.Code) {
		leader := r.MetaNodeManager.Leader()
		resp.Result.Leader = &pb.ClusterLeader{
			Available: leader != "",
			Address:   leader,
		}
	}

	return resp, nil
}

// Apply meta command
func (r *RPC) Apply(ctx context.Context, req *pb.MetaApplyReq) (*pb.MetaApplyResp, error) {
	return nil, nil
}
