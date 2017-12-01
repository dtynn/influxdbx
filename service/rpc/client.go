package rpc

import (
	"fmt"

	"github.com/dtynn/influxdbx/internal/pb"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type forwardResponse interface {
	GetResult() *pb.ClusterResult
	proto.Message
}

// MetaClient meta client interface
type MetaClient interface {
	pb.MetaClient
	WithLogger(l *zap.Logger)
}

type metaClient struct {
	target    string
	forward   bool
	maxInvoke int
	logger    *zap.Logger
}

// NewMetaClient return a meta client
func NewMetaClient(target string, forward bool) MetaClient {
	return &metaClient{
		target:    target,
		forward:   forward,
		maxInvoke: 5,
		logger:    zap.NewNop(),
	}
}

func (c *metaClient) Join(ctx context.Context, in *pb.MetaJoinReq, opts ...grpc.CallOption) (*pb.MetaJoinResp, error) {
	out := new(pb.MetaJoinResp)
	err := c.invoke(ctx, "/pb.Meta/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *metaClient) Remove(ctx context.Context, in *pb.MetaRemoveReq, opts ...grpc.CallOption) (*pb.MetaRemoveResp, error) {
	out := new(pb.MetaRemoveResp)
	err := c.invoke(ctx, "/pb.Meta/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *metaClient) Apply(ctx context.Context, in *pb.MetaApplyReq, opts ...grpc.CallOption) (*pb.MetaApplyResp, error) {
	out := new(pb.MetaApplyResp)
	err := c.invoke(ctx, "/pb.Meta/Apply", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// WithLogger setup logger
func (c *metaClient) WithLogger(l *zap.Logger) {
	c.logger = l.With(zap.String("client", "meta"))
}

func (c *metaClient) invoke(ctx context.Context, method string, in, out interface{}, opts ...grpc.CallOption) error {
	for tried := 0; tried < c.maxInvoke; tried++ {
		if tried > 0 {
			c.logger.Info(fmt.Sprintf("[%s] request is forwarded to %s, tried %d time(s).", method, c.target, tried))
		}

		// network or service error
		err := invoke(ctx, method, in, out, c.target, opts...)
		if err != nil {
			return err
		}

		if !c.forward {
			return nil
		}

		rc, ok := out.(forwardResponse)
		if !ok {
			return nil
		}

		switch rc.GetResult().GetCode() {
		case pb.ResultCode_ResultCodeErrNotLeader,
			pb.ResultCode_ResultCodeErrLeadershipLost:
			if !rc.GetResult().GetLeader().GetAvailable() {
				return nil
			}

		default:
			return nil
		}

		c.target = rc.GetResult().GetLeader().GetAddress()
		rc.Reset()
	}

	return nil
}

func invoke(ctx context.Context, method string, in, out interface{}, target string, opts ...grpc.CallOption) error {
	cc, err := defaultConnMgr.Get(target)
	if err != nil {
		return err
	}

	return cc.Invoke(ctx, method, in, out, opts...)
}
