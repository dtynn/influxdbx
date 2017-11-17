package raft

import (
	"time"

	"github.com/dtynn/influxdbx/service/cluster/proto"
	"github.com/hashicorp/raft"
	"golang.org/x/net/context"
)

var (
	joinOKResp = &proto.ClusterJoinResp{
		Code: proto.Code_CodeOK,
	}

	removeOKResp = &proto.ClusterRemoveResp{
		Code: proto.Code_CodeOK,
	}
)

func (r *Raft) Join(ctx context.Context, req *proto.ClusterJoinReq) (*proto.ClusterJoinResp, error) {
	if r.r == nil {
		return &proto.ClusterJoinResp{
			Code: proto.Code_CodeOK,
		}, nil
	}

	if r.r.State() != raft.Leader {
		return &proto.ClusterJoinResp{
			Code: proto.Code_CodeNotLeader,
			Leader: &proto.Node{
				Address: string(r.r.Leader()),
			},
		}, nil
	}

	node := req.GetNode()
	for _, one := range node {
		f := r.r.AddVoter(raft.ServerID(one.GetId()), raft.ServerAddress(one.GetAddress()), 0, 10*time.Second)
		if err := f.Error(); err != nil {
			if err == raft.ErrNotLeader {
				return &proto.ClusterJoinResp{
					Code: proto.Code_CodeNotLeader,
					Leader: &proto.Node{
						Address: string(r.r.Leader()),
					},
				}, nil
			}

			return nil, err
		}
	}

	return joinOKResp, nil
}

func (r *Raft) Remove(ctx context.Context, req *proto.ClusterRemoveReq) (*proto.ClusterRemoveResp, error) {
	if r.r == nil {
		return &proto.ClusterRemoveResp{
			Code: proto.Code_CodeOK,
		}, nil
	}

	if r.r.State() != raft.Leader {
		return &proto.ClusterRemoveResp{
			Code: proto.Code_CodeNotLeader,
			Leader: &proto.Node{
				Address: string(r.r.Leader()),
			},
		}, nil
	}

	node := req.GetNode()
	for _, one := range node {
		f := r.r.RemoveServer(raft.ServerID(one.GetId()), 0, 10*time.Second)
		if err := f.Error(); err != nil {
			if err == raft.ErrNotLeader {
				return &proto.ClusterRemoveResp{
					Code: proto.Code_CodeNotLeader,
					Leader: &proto.Node{
						Address: string(r.r.Leader()),
					},
				}, nil
			}

			return nil, err
		}
	}

	return removeOKResp, nil
}
