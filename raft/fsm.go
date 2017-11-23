package raft

import (
	"fmt"
	"io"

	"github.com/dtynn/influxdbx/raft/internal"
	"github.com/hashicorp/raft"
	"github.com/influxdata/influxql"
)

func newFSM(r *Raft) *FSM {
	return &FSM{
		r: r,
	}
}

type FSM struct {
	r *Raft
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	cmdType, id, buf := fsmCmdRead(l.Data)
	if id != nil {
		queryID := *id
		if !f.r.queryMgr.Check(queryID) {
			return newFsmCmdResponse(nil, nil)
		}

		defer f.r.queryMgr.Done(queryID)
	}

	switch cmdType {
	// continuous query
	case fsmCmdTypeMetaCreateContinuousQuery:
		return f.applyCreateContinuousQuery(buf)

	case fsmCmdTypeMetaDropContinuousQuery:
		return f.applyDropContinuousQuery(buf)

	// database
	case fsmCmdTypeMetaCreateDatabase:
		return f.applyCreateDatabase(buf)

	case fsmCmdTypeMetaCreateDatabaseWithRetentionPolicy:
		return f.applyCreateDatabaseWithRetentionPolicy(buf)

	case fsmCmdTypeMetaDropDatabase:
		return f.applyDropDatabase(buf)

	// retetion policy
	case fsmCmdTypeMetaCreateRetentionPolicy:
		return f.applyCreateRetentionPolicy(buf)

	case fsmCmdTypeMetaDropRetentionPolicy:
		return f.applyDropRetentionPolicy(buf)

	case fsmCmdTypeMetaUpdateRetentionPolicy:
		return f.applyUpdateRetentionPolicy(buf)

	// user
	case fsmCmdTypeMetaCreateUser:
		return f.applyCreateUser(buf)

	case fsmCmdTypeMetaDropUser:
		return f.applyDropUser(buf)

	case fsmCmdTypeMetaUpdateUser:
		return f.applyUpdateUser(buf)

	// privilege
	case fsmCmdTypeMetaSetAdminPrivilege:
		return f.applySetAdminPrivilege(buf)

	case fsmCmdTypeMetaSetPrivilege:
		return f.applySetPrivilege(buf)

	default:
		return newFsmCmdResponse(nil, fmt.Errorf("unknown cmd type %v", cmdType))
	}
}

func (f *FSM) applyAcquireLease(buf []byte) fsmCmdResponse {
	var cmd internal.AcquireLeaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.AcquireLease(cmd.GetName()))
}

// database
func (f *FSM) applyDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.DatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.Database(cmd.GetName()), nil)
}

func (f *FSM) applyDatabases(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.Databases(), nil)
}

func (f *FSM) applyCreateDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.CreateDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateDatabase(cmd.GetName()))
}

func (f *FSM) applyCreateDatabaseWithRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.CreateDatabaseWithRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateDatabaseWithRetentionPolicy(cmd.GetName(), retentionPolicySpecProto2Meta(cmd.GetSpec())))
}

func (f *FSM) applyDropDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.DropDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DropDatabase(cmd.GetName()))
}

// retention policy
func (f *FSM) applyCreateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.CreateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateRetentionPolicy(cmd.GetDatabase(), retentionPolicySpecProto2Meta(cmd.GetSpec()), cmd.GetMakeDefault()))
}

func (f *FSM) applyRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.RetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.RetentionPolicy(cmd.GetDatabase(), cmd.GetName()))
}

func (f *FSM) applyDropRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.DropRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DropRetentionPolicy(cmd.GetDatabase(), cmd.GetName()))
}

func (f *FSM) applyUpdateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.UpdateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.UpdateRetentionPolicy(cmd.GetDatabase(), cmd.GetName(), retentionPolicyUpdateProto2Meta(cmd.GetUpdate()), cmd.GetMakeDefault()))
}

// user manage
func (f *FSM) applyUsers(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.Users(), nil)
}

func (f *FSM) applyUser(buf []byte) fsmCmdResponse {
	var cmd internal.UserCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.User(cmd.GetName()))
}

func (f *FSM) applyCreateUser(buf []byte) fsmCmdResponse {
	var cmd internal.CreateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	res, err := f.r.MetaClient.CreateUser(cmd.GetName(), cmd.GetPassword(), cmd.GetAdmin())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyUpdateUser(buf []byte) fsmCmdResponse {
	var cmd internal.UpdateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.UpdateUser(cmd.GetName(), cmd.GetPassword())

	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropUser(buf []byte) fsmCmdResponse {
	var cmd internal.DropUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropUser(cmd.GetName())

	return newFsmCmdResponse(nil, err)
}

// user privilege
func (f *FSM) applySetAdminPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.SetAdminPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.SetAdminPrivilege(cmd.GetUsername(), cmd.GetAdmin())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applySetPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.SetPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.SetPrivilege(cmd.GetUsername(), cmd.GetDatabase(), influxql.Privilege(cmd.GetPrivilege()))
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyUserPrivileges(buf []byte) fsmCmdResponse {
	var cmd internal.UserPrivilegesCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.UserPrivileges(cmd.GetUsername()))
}

func (f *FSM) applyUserPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.UserPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.UserPrivilege(cmd.GetUsername(), cmd.GetDatabase()))
}

func (f *FSM) applyAdminUserExists(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.AdminUserExists(), nil)
}

func (f *FSM) applyAuthenticate(buf []byte) fsmCmdResponse {
	var cmd internal.AuthenticateCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.Authenticate(cmd.GetUsername(), cmd.GetPassword()))
}

func (f *FSM) applyUserCount(buf []byte) fsmCmdResponse {

	return newFsmCmdResponse(f.r.MetaClient.UserCount(), nil)
}

// shard
func (f *FSM) applyShardIDs(buf []byte) fsmCmdResponse {

	return newFsmCmdResponse(f.r.MetaClient.ShardIDs(), nil)
}

func (f *FSM) applyShardGroupsByTimeRange(buf []byte) fsmCmdResponse {
	var cmd internal.ShardGroupsByTimeRangeCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.ShardGroupsByTimeRange(cmd.GetDatabase(), cmd.GetPolicy(), nano2time(cmd.GetTmin()), nano2time(cmd.GetTmax())))
}

func (f *FSM) applyShardsByTimeRange(buf []byte) fsmCmdResponse {
	var cmd internal.ShardsByTimeRangeCmd
	protoCmdUnmarshal(buf, &cmd)

	sources := make(influxql.Sources, 0, 10)
	if err := sources.UnmarshalBinary(cmd.GetSourceData()); err != nil {
		return newFsmCmdResponse(nil, err)
	}

	return newFsmCmdResponse(f.r.MetaClient.ShardsByTimeRange(sources, nano2time(cmd.GetTmin()), nano2time(cmd.GetTmax())))
}

func (f *FSM) applyDropShard(buf []byte) fsmCmdResponse {
	var cmd internal.DropShardCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropShard(cmd.GetId())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyPruneShardGroups(buf []byte) fsmCmdResponse {

	return newFsmCmdResponse(nil, f.r.MetaClient.PruneShardGroups())

}

func (f *FSM) applyCreateShardGroup(buf []byte) fsmCmdResponse {
	var cmd internal.CreateShardGroupCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateShardGroup(cmd.GetDatabase(), cmd.GetPolicy(), nano2time(cmd.GetTimestamp())))

}

func (f *FSM) applyDeleteShardGroup(buf []byte) fsmCmdResponse {
	var cmd internal.DeleteShardGroupCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DeleteShardGroup(cmd.GetDatabase(), cmd.GetPolicy(), cmd.GetId()))

}

func (f *FSM) applyPrecreateShardGroups(buf []byte) fsmCmdResponse {
	var cmd internal.PrecreateShardGroupsCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.PrecreateShardGroups(nano2time(cmd.GetFrom()), nano2time(cmd.GetTo())))

}

func (f *FSM) applyShardOwner(buf []byte) fsmCmdResponse {
	var cmd internal.ShardOwnerCmd
	protoCmdUnmarshal(buf, &cmd)

	var sor shardOwnerResult

	sor.database, sor.policy, sor.info = f.r.MetaClient.ShardOwner(cmd.GetId())

	return newFsmCmdResponse(sor, nil)

}

// continuous query
func (f *FSM) applyCreateContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.CreateContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.CreateContinuousQuery(cmd.GetDatabase(), cmd.GetName(), cmd.GetQuery())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.DropContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropContinuousQuery(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

// subscription
func (f *FSM) applyCreateSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.CreateSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.CreateSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName(), cmd.GetMode(), cmd.GetDestinations())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDropSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.DropSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{}, nil
}

func (f *FSM) Restore(r io.ReadCloser) error {
	return fmt.Errorf("not implemented")
}

type FSMSnapshot struct {
}

func (f *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return fmt.Errorf("not implemented")
}

func (f *FSMSnapshot) Release() {

}
