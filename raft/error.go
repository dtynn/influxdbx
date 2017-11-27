package raft

import (
	"github.com/hashicorp/raft"
)

func IsErrNotLeader(err error) bool {
	return err == ErrNotLeader
}

var ErrNotLeader = raft.ErrNotLeader
