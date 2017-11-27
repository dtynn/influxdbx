package cluster

import (
	"fmt"

	"github.com/dtynn/influxdbx/cluster/proto"
	"github.com/hashicorp/raft"
)

func Error(code proto.Code) error {
	switch code {
	case proto.Code_CodeOK:
		return nil

	case proto.Code_CodeNotLeader:
		return raft.ErrNotLeader

	case proto.Code_CodeErrLeadershipLost:
		return raft.ErrLeadershipLost

	default:
		return fmt.Errorf("unknown error code %d", code)
	}
}

func Code(err error) proto.Code {
	switch err {
	case nil:
		return proto.Code_CodeOK

	case raft.ErrNotLeader:
		return proto.Code_CodeNotLeader

	case raft.ErrLeadershipLost:
		return proto.Code_CodeErrLeadershipLost

	default:
		return proto.Code_CodeUnknown
	}
}
