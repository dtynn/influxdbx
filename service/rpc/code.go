package rpc

import (
	"errors"

	"github.com/dtynn/influxdbx/internal/pb"
	"github.com/hashicorp/raft"
)

func canForward(code pb.ResultCode) bool {
	return code == pb.ResultCode_ResultCodeErrNotLeader || code == pb.ResultCode_ResultCodeErrLeadershipLost
}

// Code2Err convert result code and msg to local error
func Code2Err(code pb.ResultCode, msg string) error {
	switch code {
	case pb.ResultCode_ResultCodeOK:
		return nil

	case pb.ResultCode_ResultCodeErrNotLeader:
		return raft.ErrNotLeader

	case pb.ResultCode_ResultCodeErrLeadershipLost:
		return raft.ErrLeadershipLost

	default:
		return errors.New(msg)
	}
}

// Err2Code convert local error to result code
func Err2Code(err error) (pb.ResultCode, string) {
	if err == nil {
		return pb.ResultCode_ResultCodeOK, ""
	}

	switch err {
	case raft.ErrNotLeader:
		return pb.ResultCode_ResultCodeErrNotLeader, err.Error()

	case raft.ErrLeadershipLost:
		return pb.ResultCode_ResultCodeErrLeadershipLost, err.Error()

	default:
		return pb.ResultCode_ResultCodeUnknown, err.Error()
	}
}
