package cluster

import (
	"github.com/dtynn/influxdbx/cluster/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type redirectResponse interface {
	GetResult() *proto.Result
	GetLeader() *proto.Leader
	Reset()
}

type clusterClient struct {
	target    string
	relay     bool
	maxInvoke int
}

func NewClusterClient(target string, relay bool) proto.ClusterClient {
	return &clusterClient{
		target:    target,
		relay:     relay,
		maxInvoke: 5,
	}
}

func (c *clusterClient) Join(ctx context.Context, in *proto.ClusterJoinReq, opts ...grpc.CallOption) (*proto.ClusterJoinResp, error) {
	out := new(proto.ClusterJoinResp)
	err := c.invoke(ctx, "/proto.Cluster/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Remove(ctx context.Context, in *proto.ClusterRemoveReq, opts ...grpc.CallOption) (*proto.ClusterRemoveResp, error) {
	out := new(proto.ClusterRemoveResp)
	err := c.invoke(ctx, "/proto.Cluster/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *clusterClient) invoke(ctx context.Context, method string, in, out interface{}, opts ...grpc.CallOption) error {
	for tried := 0; tried < c.maxInvoke; tried++ {
		// network or service error
		err := invoke(ctx, method, in, out, c.target, opts...)
		if err != nil {
			return err
		}

		if !c.relay {
			return nil
		}

		rc, ok := out.(redirectResponse)
		if !ok {
			return nil
		}

		if rc.GetResult().GetCode() != proto.Code_CodeNotLeader || !rc.GetLeader().GetAvailable() {
			return nil
		}

		c.target = rc.GetLeader().GetAddress()
		rc.Reset()
	}

	return nil
}

func invoke(ctx context.Context, method string, in, out interface{}, target string, opts ...grpc.CallOption) error {
	cc, err := getConn(target)
	if err != nil {
		return err
	}

	return cc.Invoke(ctx, method, in, out, opts...)
}
