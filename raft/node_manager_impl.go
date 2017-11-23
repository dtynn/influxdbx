package raft

import (
	"time"

	"github.com/hashicorp/raft"
)

var (
	nodeOperationTimeout = 5 * time.Second
)

func (r *Raft) AddNode(id, address string) error {
	return r.r.AddVoter(raft.ServerID(id), raft.ServerAddress(address), 0, nodeOperationTimeout).Error()
}

func (r *Raft) RemoveNode(id string) error {
	return r.r.RemoveServer(raft.ServerID(id), 0, nodeOperationTimeout).Error()
}

func (r *Raft) Leader() (string, bool) {
	addr := string(r.r.Leader())
	return addr, addr != ""
}

func (r *Raft) IsLeader() bool {
	return r.r != nil && r.r.State() == raft.Leader
}

func (r *Raft) IsErrNotLeader(err error) bool {
	return err == raft.ErrNotLeader
}
