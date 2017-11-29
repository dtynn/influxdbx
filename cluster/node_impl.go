package cluster

import (
	"github.com/dtynn/influxdbx/cluster/proto"
	"golang.org/x/net/context"
)

var (
	joinOKResp = &proto.ClusterJoinResp{
		Result: &proto.Result{
			Code: proto.Code_CodeOK,
		},
	}

	removeOKResp = &proto.ClusterRemoveResp{
		Result: &proto.Result{
			Code: proto.Code_CodeOK,
		},
	}
)

func (c *Cluster) Join(ctx context.Context, req *proto.ClusterJoinReq) (*proto.ClusterJoinResp, error) {
	if !c.NodeManager.IsLeader() {
		leader := &proto.Leader{}
		leader.Address, leader.Available = c.NodeManager.Leader()

		return &proto.ClusterJoinResp{
			Result: &proto.Result{
				Code: proto.Code_CodeNotLeader,
			},
			Leader: leader,
		}, nil
	}

	for _, one := range req.GetNode() {
		err := c.NodeManager.AddNode(one.GetId(), one.GetAddress())
		if err == nil {
			continue
		}

		if c.NodeManager.IsErrNotLeader(err) {
			leader := &proto.Leader{}
			leader.Address, leader.Available = c.NodeManager.Leader()

			return &proto.ClusterJoinResp{
				Result: &proto.Result{
					Code: proto.Code_CodeNotLeader,
				},
				Leader: leader,
			}, nil
		}

		return nil, err
	}

	return joinOKResp, nil
}

func (c *Cluster) Remove(ctx context.Context, req *proto.ClusterRemoveReq) (*proto.ClusterRemoveResp, error) {
	if !c.NodeManager.IsLeader() {
		leader := &proto.Leader{}
		leader.Address, leader.Available = c.NodeManager.Leader()

		return &proto.ClusterRemoveResp{
			Result: &proto.Result{
				Code: proto.Code_CodeNotLeader,
			},
			Leader: leader,
		}, nil
	}

	for _, one := range req.GetNode() {
		err := c.NodeManager.RemoveNode(one.GetId())
		if err == nil {
			continue
		}

		if c.NodeManager.IsErrNotLeader(err) {
			leader := &proto.Leader{}
			leader.Address, leader.Available = c.NodeManager.Leader()

			return &proto.ClusterRemoveResp{
				Result: &proto.Result{
					Code: proto.Code_CodeNotLeader,
				},
				Leader: leader,
			}, nil
		}

		return nil, err
	}

	return removeOKResp, nil
}
