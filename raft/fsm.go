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
	// database
	case fsmCmdTypeMetaCreateDatabase:
		return f.applyMetaCreateDatabase(buf)

	case fsmCmdTypeMetaCreateDatabaseWithRetentionPolicy:
		return f.applyMetaCreateDatabaseWithRetentionPolicy(buf)

	case fsmCmdTypeMetaDropDatabase:
		return f.applyMetaDropDatabase(buf)

	// retetion policy
	case fsmCmdTypeMetaCreateRetentionPolicy:
		return f.applyMetaCreateRetentionPolicy(buf)

	case fsmCmdTypeMetaDropRetentionPolicy:
		return f.applyMetaDropRetentionPolicy(buf)

	case fsmCmdTypeMetaUpdateRetentionPolicy:
		return f.applyMetaUpdateRetentionPolicy(buf)

	// user
	case fsmCmdTypeMetaCreateUser:
		return f.applyMetaCreateUser(buf)

	case fsmCmdTypeMetaDropUser:
		return f.applyMetaDropUser(buf)

	case fsmCmdTypeMetaUpdateUser:
		return f.applyMetaUpdateUser(buf)

	// privilege
	case fsmCmdTypeMetaSetAdminPrivilege:
		return f.applyMetaSetAdminPrivilege(buf)

	case fsmCmdTypeMetaSetPrivilege:
		return f.applyMetaSetPrivilege(buf)

	// shad
	case fsmCmdTypeMetaDropShard:
		return f.applyMetaDropShard(buf)

	case fsmCmdTypeMetaPruneShardGroups:
		return f.applyMetaPruneShardGroups(buf)

	case fsmCmdTypeMetaCreateShardGroup:
		return f.applyMetaCreateShardGroup(buf)

	case fsmCmdTypeMetaDeleteShardGroup:
		return f.applyMetaDeleteShardGroup(buf)

	case fsmCmdTypeMetaPrecreateShardGroups:
		return f.applyMetaPrecreateShardGroups(buf)

	// continuous query
	case fsmCmdTypeMetaCreateContinuousQuery:
		return f.applyMetaCreateContinuousQuery(buf)

	case fsmCmdTypeMetaDropContinuousQuery:
		return f.applyMetaDropContinuousQuery(buf)

	case fsmCmdTypeMetaCreateSubscription:
		return f.applyMetaCreateSubscription(buf)

	case fsmCmdTypeMetaDropSubscription:
		return f.applyMetaDropSubscription(buf)

	case fsmCmdTypeMetaAcquireLease:
		return f.applyMetaAcquireLease(buf)

	case fsmCmdTypeMetaDatabase:
		return f.applyMetaDatabase(buf)

	case fsmCmdTypeMetaDatabases:
		return f.applyMetaDatabases(buf)

	case fsmCmdTypeMetaRetentionPolicy:
		return f.applyMetaRetentionPolicy(buf)

	case fsmCmdTypeMetaUsers:
		return f.applyMetaUsers(buf)

	case fsmCmdTypeMetaUser:
		return f.applyMetaUser(buf)

	case fsmCmdTypeMetaUserPrivileges:
		return f.applyMetaUserPrivileges(buf)

	case fsmCmdTypeMetaUserPrivilege:
		return f.applyMetaUserPrivilege(buf)

	case fsmCmdTypeMetaAdminUserExists:
		return f.applyMetaAdminUserExists(buf)

	case fsmCmdTypeMetaAuthenticate:
		return f.applyMetaAuthenticate(buf)

	case fsmCmdTypeMetaShardGroupsByTimeRange:
		return f.applyMetaShardGroupsByTimeRange(buf)

	case fsmCmdTypeDataCreateShard:
		return f.applyDataCreateShard(buf)

	case fsmCmdTypeDataWriteToShard:
		return f.applyDataWriteToShard(buf)

	case fsmCmdTypeDataDeleteShard:
		return f.applyDataDeleteShard(buf)

	case fsmCmdTypeDataShardGroup:
		return f.applyDataShardGroup(buf)

	case fsmCmdTypeDataShardIDs:
		return f.applyDataShardIDs(buf)

	case fsmCmdTypeDataDeleteDatabase:
		return f.applyDataDeleteDatabase(buf)

	case fsmCmdTypeDataDeleteRetentionPolicy:
		return f.applyDataDeleteRetentionPolicy(buf)

	case fsmCmdTypeDataDeleteMeasurement:
		return f.applyDataDeleteMeasurement(buf)

	case fsmCmdTypeDataMeasurementNames:
		return f.applyDataMeasurementNames(buf)

	case fsmCmdTypeDataMeasurementsCardinality:
		return f.applyDataMeasurementsCardinality(buf)

	case fsmCmdTypeDataDeleteSeries:
		return f.applyDataDeleteSeries(buf)

	case fsmCmdTypeDataSeriesCardinality:
		return f.applyDataSeriesCardinality(buf)

	case fsmCmdTypeDataTagKeys:
		return f.applyDataTagKeys(buf)

	case fsmCmdTypeDataTagValues:
		return f.applyDataTagValues(buf)

	default:
		return newFsmCmdResponse(nil, fmt.Errorf("unknown cmd type %v", cmdType))
	}
}

func (f *FSM) applyMetaAcquireLease(buf []byte) fsmCmdResponse {
	var cmd internal.MetaAcquireLeaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.AcquireLease(cmd.GetName()))
}

// database
func (f *FSM) applyMetaDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.Database(cmd.GetName()), nil)
}

func (f *FSM) applyMetaDatabases(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.Databases(), nil)
}

func (f *FSM) applyMetaCreateDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateDatabase(cmd.GetName()))
}

func (f *FSM) applyMetaCreateDatabaseWithRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateDatabaseWithRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateDatabaseWithRetentionPolicy(cmd.GetName(), retentionPolicySpecProto2Meta(cmd.GetSpec())))
}

func (f *FSM) applyMetaDropDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DropDatabase(cmd.GetName()))
}

// retention policy
func (f *FSM) applyMetaCreateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateRetentionPolicy(cmd.GetDatabase(), retentionPolicySpecProto2Meta(cmd.GetSpec()), cmd.GetMakeDefault()))
}

func (f *FSM) applyMetaRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.MetaRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.RetentionPolicy(cmd.GetDatabase(), cmd.GetName()))
}

func (f *FSM) applyMetaDropRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DropRetentionPolicy(cmd.GetDatabase(), cmd.GetName()))
}

func (f *FSM) applyMetaUpdateRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.MetaUpdateRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.UpdateRetentionPolicy(cmd.GetDatabase(), cmd.GetName(), retentionPolicyUpdateProto2Meta(cmd.GetUpdate()), cmd.GetMakeDefault()))
}

// user manage
func (f *FSM) applyMetaUsers(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.Users(), nil)
}

func (f *FSM) applyMetaUser(buf []byte) fsmCmdResponse {
	var cmd internal.MetaUserCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.User(cmd.GetName()))
}

func (f *FSM) applyMetaCreateUser(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	res, err := f.r.MetaClient.CreateUser(cmd.GetName(), cmd.GetPassword(), cmd.GetAdmin())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyMetaUpdateUser(buf []byte) fsmCmdResponse {
	var cmd internal.MetaUpdateUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.UpdateUser(cmd.GetName(), cmd.GetPassword())

	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaDropUser(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropUserCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropUser(cmd.GetName())

	return newFsmCmdResponse(nil, err)
}

// user privilege
func (f *FSM) applyMetaSetAdminPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.MetaSetAdminPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.SetAdminPrivilege(cmd.GetUsername(), cmd.GetAdmin())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaSetPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.MetaSetPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.SetPrivilege(cmd.GetUsername(), cmd.GetDatabase(), influxql.Privilege(cmd.GetPrivilege()))
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaUserPrivileges(buf []byte) fsmCmdResponse {
	var cmd internal.MetaUserPrivilegesCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.UserPrivileges(cmd.GetUsername()))
}

func (f *FSM) applyMetaUserPrivilege(buf []byte) fsmCmdResponse {
	var cmd internal.MetaUserPrivilegeCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.UserPrivilege(cmd.GetUsername(), cmd.GetDatabase()))
}

func (f *FSM) applyMetaAdminUserExists(buf []byte) fsmCmdResponse {
	return newFsmCmdResponse(f.r.MetaClient.AdminUserExists(), nil)
}

func (f *FSM) applyMetaAuthenticate(buf []byte) fsmCmdResponse {
	var cmd internal.MetaAuthenticateCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.Authenticate(cmd.GetUsername(), cmd.GetPassword()))
}

// shard
func (f *FSM) applyMetaShardGroupsByTimeRange(buf []byte) fsmCmdResponse {
	var cmd internal.MetaShardGroupsByTimeRangeCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.ShardGroupsByTimeRange(cmd.GetDatabase(), cmd.GetPolicy(), nano2time(cmd.GetTmin()), nano2time(cmd.GetTmax())))
}

func (f *FSM) applyMetaDropShard(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropShardCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropShard(cmd.GetId())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaPruneShardGroups(buf []byte) fsmCmdResponse {

	return newFsmCmdResponse(nil, f.r.MetaClient.PruneShardGroups())

}

func (f *FSM) applyMetaCreateShardGroup(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateShardGroupCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(f.r.MetaClient.CreateShardGroup(cmd.GetDatabase(), cmd.GetPolicy(), nano2time(cmd.GetTimestamp())))

}

func (f *FSM) applyMetaDeleteShardGroup(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDeleteShardGroupCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.DeleteShardGroup(cmd.GetDatabase(), cmd.GetPolicy(), cmd.GetId()))

}

func (f *FSM) applyMetaPrecreateShardGroups(buf []byte) fsmCmdResponse {
	var cmd internal.MetaPrecreateShardGroupsCmd
	protoCmdUnmarshal(buf, &cmd)

	return newFsmCmdResponse(nil, f.r.MetaClient.PrecreateShardGroups(nano2time(cmd.GetFrom()), nano2time(cmd.GetTo())))

}

// continuous query
func (f *FSM) applyMetaCreateContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.CreateContinuousQuery(cmd.GetDatabase(), cmd.GetName(), cmd.GetQuery())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaDropContinuousQuery(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropContinuousQueryCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropContinuousQuery(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

// subscription
func (f *FSM) applyMetaCreateSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.MetaCreateSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.CreateSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName(), cmd.GetMode(), cmd.GetDestinations())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyMetaDropSubscription(buf []byte) fsmCmdResponse {
	var cmd internal.MetaDropSubscriptionCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.MetaClient.DropSubscription(cmd.GetDatabase(), cmd.GetRp(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataCreateShard(buf []byte) fsmCmdResponse {
	var cmd internal.DataCreateShardCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.TSDBStore.CreateShard(cmd.GetDatabase(), cmd.GetPolicy(), cmd.GetShardID(), cmd.GetEnabled())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataWriteToShard(buf []byte) fsmCmdResponse {
	var cmd internal.DataWriteToShardCmd
	protoCmdUnmarshal(buf, &cmd)

	points, err := pointsProto2Data(cmd.GetPointData())
	if err != nil {
		return newFsmCmdResponse(nil, err)
	}

	err = f.r.TSDBStore.WriteToShard(cmd.GetShardID(), points)
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataDeleteShard(buf []byte) fsmCmdResponse {
	var cmd internal.DataDeleteShardCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.TSDBStore.DeleteShard(cmd.GetId())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataShardGroup(buf []byte) fsmCmdResponse {
	var cmd internal.DataShardGroupCmd
	protoCmdUnmarshal(buf, &cmd)

	group := f.r.TSDBStore.ShardGroup(cmd.GetIds())
	return newFsmCmdResponse(group, nil)
}

func (f *FSM) applyDataShardIDs(buf []byte) fsmCmdResponse {

	ids := f.r.TSDBStore.ShardIDs()
	return newFsmCmdResponse(ids, nil)
}

func (f *FSM) applyDataDeleteDatabase(buf []byte) fsmCmdResponse {
	var cmd internal.DataDeleteDatabaseCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.TSDBStore.DeleteDatabase(cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataDeleteRetentionPolicy(buf []byte) fsmCmdResponse {
	var cmd internal.DataDeleteRetentionPolicyCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.TSDBStore.DeleteRetentionPolicy(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataDeleteMeasurement(buf []byte) fsmCmdResponse {
	var cmd internal.DataDeleteMeasurementCmd
	protoCmdUnmarshal(buf, &cmd)

	err := f.r.TSDBStore.DeleteMeasurement(cmd.GetDatabase(), cmd.GetName())
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataMeasurementNames(buf []byte) fsmCmdResponse {
	var cmd internal.DataMeasurementNamesCmd
	protoCmdUnmarshal(buf, &cmd)

	cond, err := influxql.ParseExpr(cmd.GetCondition())
	if err != nil {
		return newFsmCmdResponse(nil, err)
	}

	res, err := f.r.TSDBStore.MeasurementNames(protoUser2Authorizer(cmd.GetUser()), cmd.GetDatabase(), cond)
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyDataMeasurementsCardinality(buf []byte) fsmCmdResponse {
	var cmd internal.DataMeasurementsCardinalityCmd
	protoCmdUnmarshal(buf, &cmd)

	res, err := f.r.TSDBStore.MeasurementsCardinality(cmd.GetDatabase())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyDataDeleteSeries(buf []byte) fsmCmdResponse {
	var cmd internal.DataDeleteSeriesCmd
	protoCmdUnmarshal(buf, &cmd)

	sources := make(influxql.Sources, 0, 10)
	if err := sources.UnmarshalBinary(cmd.GetSourceData()); err != nil {
		return newFsmCmdResponse(nil, err)
	}

	condition, err := influxql.ParseExpr(cmd.GetCondition())
	if err != nil {
		return newFsmCmdResponse(nil, err)
	}

	err = f.r.TSDBStore.DeleteSeries(cmd.GetDatabase(), []influxql.Source(sources), condition)
	return newFsmCmdResponse(nil, err)
}

func (f *FSM) applyDataSeriesCardinality(buf []byte) fsmCmdResponse {
	var cmd internal.DataSeriesCardinalityCmd
	protoCmdUnmarshal(buf, &cmd)

	res, err := f.r.TSDBStore.SeriesCardinality(cmd.GetDatabase())
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyDataTagKeys(buf []byte) fsmCmdResponse {
	var cmd internal.DataTagKeysCmd
	protoCmdUnmarshal(buf, &cmd)

	cond, err := influxql.ParseExpr(cmd.GetCondition())
	if err != nil {
		return newFsmCmdResponse(nil, err)
	}

	res, err := f.r.TSDBStore.TagKeys(protoUser2Authorizer(cmd.GetUser()), cmd.GetShardIDs(), cond)
	return newFsmCmdResponse(res, err)
}

func (f *FSM) applyDataTagValues(buf []byte) fsmCmdResponse {
	var cmd internal.DataTagValuesCmd
	protoCmdUnmarshal(buf, &cmd)

	cond, err := influxql.ParseExpr(cmd.GetCondition())
	if err != nil {
		return newFsmCmdResponse(nil, err)
	}

	res, err := f.r.TSDBStore.TagValues(protoUser2Authorizer(cmd.GetUser()), cmd.GetShardIDs(), cond)
	return newFsmCmdResponse(res, err)
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{}, nil
}

func (f *FSM) Restore(r io.ReadCloser) error {
	return fmt.Errorf("FSM.Resotre not implemented")
}

type FSMSnapshot struct {
}

func (f *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return fmt.Errorf("FSMSnapshot.Persist not implemented")
}

func (f *FSMSnapshot) Release() {

}
