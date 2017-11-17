package raft

import (
	"fmt"
	"io"

	"github.com/dtynn/influxdbx/service/raft/internal"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	"github.com/influxdata/influxdb/coordinator"
)

type FSM struct {
	coordinator.MetaClient
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	cmdType := fsmCmdType(l.Data[0])
	buf := l.Data[1:]

	switch cmdType {
	case fsmCmdTypeCreateUser:
		return f.applyCreateUser(buf)

	case fsmCmdTypeUpdateUser:
		return f.applyUpdateUser(buf)

	case fsmCmdTypeDropUser:
		return f.applyDropUser(buf)

	case fsmCmdTypeUsers:
		return f.applyUsers(buf)

	case fsmCmdTypeUserPrivilege:
		return f.applyUserPrivilege(buf)

	case fsmCmdTypeUserPrivileges:
		return f.applyUserPrivileges(buf)

	default:
		return newFsmCmdResponse(nil, fmt.Errorf("unknown cmd type %v", cmdType))
	}
}

func (f *FSM) applyCreateUser(buf []byte) fsmCmdResponse {
	var cmd internal.CreateUserCmd
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		panic(fmt.Errorf("malformed create user cmd %s", err))
	}

	res, err := f.MetaClient.CreateUser(cmd.GetName(), cmd.GetPassword(), cmd.GetAdmin())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyUpdateUser(buf []byte) fsmCmdResponse {
	var cmd internal.UpdateUserCmd
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		panic(fmt.Errorf("malformed update user cmd %s", err))
	}

	err := f.MetaClient.UpdateUser(cmd.GetName(), cmd.GetPassword())

	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropUser(buf []byte) fsmCmdResponse {
	var cmd internal.DropUserCmd
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		panic(fmt.Errorf("malformed drop user cmd %s", err))
	}

	err := f.MetaClient.DropUser(cmd.GetName())

	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyUsers(buf []byte) fsmCmdResponse {
	res := f.MetaClient.Users()
	return newFsmCmdResponse(res, nil)
}

func (f *FSM) applyUserPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.UserPrivilegeCmd
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		panic(fmt.Errorf("malformed user privilege cmd %s", err))
	}

	res, err := f.MetaClient.UserPrivilege(cmd.GetUsername(), cmd.GetDatabase())

	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyUserPrivileges(buf []byte) fsmCmdResponse {
	var cmd internal.UserPrivilegesCmd
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		panic(fmt.Errorf("malformed user privileges cmd %s", err))
	}

	res, err := f.MetaClient.UserPrivileges(cmd.GetUsername())

	return newFsmCmdResponse(res, err)
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{}, nil
}

func (f *FSM) Restore(r io.ReadCloser) error {
	return nil
}

type FSMSnapshot struct {
}

func (f *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (f *FSMSnapshot) Release() {

}
