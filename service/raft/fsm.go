package raft

import (
	"fmt"
	"io"

	"github.com/dtynn/influxdbx/service/raft/internal"
	"github.com/hashicorp/raft"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

type FSM struct {
	MetaClient *meta.Client
	queryMgr   *localQueryManager
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	cmdType, id, buf := fsmCmdRead(l.Data)
	if id != nil {
		queryID := *id
		if !f.queryMgr.Check(queryID) {
			return newFsmCmdResponse(nil, nil)
		}

		defer f.queryMgr.Done(queryID)
	}

	switch cmdType {
	// continuous query
	case fsmCmdTypeCreateContinuousQuery:
		return f.applyCreateContinuousQuery(buf)

	case fsmCmdTypeDropContinuousQuery:
		return f.applyDropContinuousQuery(buf)

	// database
	case fsmCmdTypeCreateDatabase:
		return f.applyCreateDatabase(buf)

	case fsmCmdTypeCreateDatabaseWithRetentionPolicy:
		return f.applyCreateDatabaseWithRetentionPolicy(buf)

	case fsmCmdTypeDropDatabase:
		return f.applyDropDatabase(buf)

	// retetion policy
	case fsmCmdTypeCreateRetentionPolicy:
		return f.applyCreateRetentionPolicy(buf)

	case fsmCmdTypeDropRetentionPolicy:
		return f.applyDropRetentionPolicy(buf)

	case fsmCmdTypeUpdateRetentionPolicy:
		return f.applyUpdateRetentionPolicy(buf)

	// user
	case fsmCmdTypeCreateUser:
		return f.applyCreateUser(buf)

	case fsmCmdTypeDropUser:
		return f.applyDropUser(buf)

	case fsmCmdTypeUpdateUser:
		return f.applyUpdateUser(buf)

	// privilege
	case fsmCmdTypeSetAdminPrivilege:
		return f.applySetAdminPrivilege(buf)

	case fsmCmdTypeSetPrivilege:
		return f.applySetPrivilege(buf)

	default:
		return newFsmCmdResponse(nil, fmt.Errorf("unknown cmd type %v", cmdType))
	}
}

// continuous query
func (f *FSM) applyCreateContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.CreateContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.CreateContinuousQuery(cmd.GetDatabase(), cmd.GetName(), cmd.GetQuery())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.DropContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropContinuousQuery(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

// database
func (f *FSM) applyCreateDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.CreateDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	dbinfo, err := f.MetaClient.CreateDatabase(cmd.GetName())
	return newFsmCmdResponse(dbinfo, err)
}

func (f *FSM) applyCreateDatabaseWithRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.CreateDatabaseWithRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	dbinfo, err := f.MetaClient.CreateDatabaseWithRetentionPolicy(cmd.GetName(), retentionPolicySpecProto2Meta(cmd.GetSpec()))
	return newFsmCmdResponse(dbinfo, err)
}

func (f *FSM) applyDropDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.DropDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropDatabase(cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

// retention policy
func (f *FSM) applyCreateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.CreateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	pInfo, err := f.MetaClient.CreateRetentionPolicy(cmd.GetDatabase(), retentionPolicySpecProto2Meta(cmd.GetSpec()), cmd.GetMakeDefault())
	return newFsmCmdResponse(pInfo, err)
}

func (f *FSM) applyDropRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.DropRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropRetentionPolicy(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyUpdateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.UpdateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.UpdateRetentionPolicy(cmd.GetDatabase(), cmd.GetName(), retentionPolicyUpdateProto2Meta(cmd.GetUpdate()), cmd.GetMakeDefault())
	return newFsmCmdResponse(nil, err)
}

// subscription
func (f *FSM) applyCreateSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.CreateSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.CreateSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName(), cmd.GetMode(), cmd.GetDestinations())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.DropSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

// user manage
func (f *FSM) applyCreateUser(buf []byte) fsmCmdResponse {
	var cmd internal.CreateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	res, err := f.MetaClient.CreateUser(cmd.GetName(), cmd.GetPassword(), cmd.GetAdmin())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyDropUser(buf []byte) fsmCmdResponse {
	var cmd internal.DropUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropUser(cmd.GetName())

	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyUpdateUser(buf []byte) fsmCmdResponse {
	var cmd internal.UpdateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.UpdateUser(cmd.GetName(), cmd.GetPassword())

	return newFsmCmdResponse(nil, err)
}

// user privilege
func (f *FSM) applySetAdminPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.SetAdminPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.SetAdminPrivilege(cmd.GetUsername(), cmd.GetAdmin())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applySetPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.SetPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.SetPrivilege(cmd.GetUsername(), cmd.GetDatabase(), influxql.Privilege(cmd.GetPrivilege()))
	return newFsmCmdResponse(nil, err)
}

// shard
func (f *FSM) applyDropShard(buf []byte) fsmCmdResponse {
	var cmd internal.DropShardCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.MetaClient.DropShard(cmd.GetId())
	return newFsmCmdResponse(nil, err)
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
