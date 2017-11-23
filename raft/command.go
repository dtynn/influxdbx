package raft

import (
	"fmt"
	"time"

	"github.com/dtynn/influxdbx/raft/internal"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxdb/uuid"
)

const (
	nanoi64 = int64(time.Nanosecond)
)

type fsmCmdType byte

const (
	fsmCmdTypeUnknown fsmCmdType = 0

	fsmCmdTypeMetaCreateDatabase                    = 1
	fsmCmdTypeMetaCreateDatabaseWithRetentionPolicy = 2
	fsmCmdTypeMetaDropDatabase                      = 3
	fsmCmdTypeMetaCreateRetentionPolicy             = 4
	fsmCmdTypeMetaDropRetentionPolicy               = 5
	fsmCmdTypeMetaUpdateRetentionPolicy             = 6
	fsmCmdTypeMetaCreateUser                        = 7
	fsmCmdTypeMetaUpdateUser                        = 8
	fsmCmdTypeMetaDropUser                          = 9
	fsmCmdTypeMetaSetPrivilege                      = 10
	fsmCmdTypeMetaSetAdminPrivilege                 = 11
	fsmCmdTypeMetaDropShard                         = 12
	fsmCmdTypeMetaPruneShardGroups                  = 13
	fsmCmdTypeMetaCreateShardGroup                  = 14
	fsmCmdTypeMetaDeleteShardGroup                  = 15
	fsmCmdTypeMetaPrecreateShardGroups              = 16
	fsmCmdTypeMetaCreateContinuousQuery             = 17
	fsmCmdTypeMetaDropContinuousQuery               = 18
	fsmCmdTypeMetaCreateSubscription                = 19
	fsmCmdTypeMetaDropSubscription                  = 20

	fsmCmdTypeMetaAcquireLease           = 101
	fsmCmdTypeMetaClusterID              = 102
	fsmCmdTypeMetaDatabase               = 103
	fsmCmdTypeMetaDatabases              = 104
	fsmCmdTypeMetaRetentionPolicy        = 105
	fsmCmdTypeMetaUsers                  = 106
	fsmCmdTypeMetaUser                   = 107
	fsmCmdTypeMetaUserPrivileges         = 108
	fsmCmdTypeMetaUserPrivilege          = 109
	fsmCmdTypeMetaAdminUserExists        = 110
	fsmCmdTypeMetaAuthenticate           = 111
	fsmCmdTypeMetaUserCount              = 112
	fsmCmdTypeMetaShardIDs               = 113
	fsmCmdTypeMetaShardGroupsByTimeRange = 114
	fsmCmdTypeMetaShardsByTimeRange      = 115
	fsmCmdTypeMetaShardOwner             = 116
)

func (f fsmCmdType) isLocalQuery() bool {
	return f/100 > 0
}

type fsmCmdResponse struct {
	res interface{}
	err error
}

func newFsmCmdResponse(res interface{}, err error) fsmCmdResponse {
	return fsmCmdResponse{
		res: res,
		err: err,
	}
}

func fsmCmdMarshal(cmdType fsmCmdType, cmd proto.Message) ([]byte, *uuid.UUID, error) {
	var b []byte
	var err error

	if cmd != nil {
		b, err = proto.Marshal(cmd)
		if err != nil {
			return nil, nil, err
		}
	}

	var idPtr *uuid.UUID

	isLocalQuery := cmdType.isLocalQuery()

	size := len(b) + 1
	if isLocalQuery {
		size += 16
	}

	bs := make([]byte, size)
	bs[0] = byte(cmdType)

	head := 1
	if isLocalQuery {
		head = 17
		id := uuid.TimeUUID()
		copy(bs[1:17], id[:])
		idPtr = &id
	}

	copy(bs[head:], b)

	return bs, idPtr, nil
}

func fsmCmdRead(data []byte) (t fsmCmdType, id *uuid.UUID, buf []byte) {
	t = fsmCmdType(data[0])
	head := 1

	if t.isLocalQuery() {
		var queryID uuid.UUID
		copy(data[1:17], id[:])
		id = &queryID

		head = 17
	}

	buf = data[head:]

	return
}

func protoCmdUnmarshal(buf []byte, cmd proto.Message) {
	if err := proto.Unmarshal(buf, cmd); err != nil {
		panic(fmt.Errorf("malformed data for %T: %s", cmd, err))
	}
}

func retentionPolicySpecProto2Meta(pSpec *internal.RetentionPolicySpec) *meta.RetentionPolicySpec {
	var spec *meta.RetentionPolicySpec
	if pSpec != nil {
		spec = &meta.RetentionPolicySpec{
			Name:               pSpec.GetName(),
			ReplicaN:           pSpec.GetReplicaN().IntPtr(),
			Duration:           pSpec.GetDuration().DurationPtr(),
			ShardGroupDuration: time.Duration(pSpec.GetShardGroupDuration()),
		}
	}

	return spec
}

func retentionPolicyUpdateProto2Meta(pUpdate *internal.RetentionPolicyUpdate) *meta.RetentionPolicyUpdate {
	var update *meta.RetentionPolicyUpdate
	if pUpdate != nil {
		update = &meta.RetentionPolicyUpdate{
			Name:               pUpdate.GetName().StringPtr(),
			ReplicaN:           pUpdate.GetReplicaN().IntPtr(),
			Duration:           pUpdate.GetDuration().DurationPtr(),
			ShardGroupDuration: pUpdate.GetShardGroupDuration().DurationPtr(),
		}
	}

	return update
}

func retentionPolicySpecMeta2Proto(spec *meta.RetentionPolicySpec) *internal.RetentionPolicySpec {
	var pSpec *internal.RetentionPolicySpec
	if spec != nil {
		pSpec = &internal.RetentionPolicySpec{
			Name:               spec.Name,
			ShardGroupDuration: int64(spec.ShardGroupDuration),
		}

		if spec.ReplicaN != nil {
			pSpec.ReplicaN = &internal.OptionalInt64{
				Val: int64(*spec.ReplicaN),
			}
		}

		if spec.Duration != nil {
			pSpec.Duration = &internal.OptionalInt64{
				Val: int64(*spec.Duration),
			}
		}
	}

	return pSpec
}

func retentionPolicyUpdateMeta2Proto(update *meta.RetentionPolicyUpdate) *internal.RetentionPolicyUpdate {
	var pUpdate *internal.RetentionPolicyUpdate
	if update != nil {
		pUpdate = &internal.RetentionPolicyUpdate{}

		if update.Name != nil {
			pUpdate.Name = &internal.OptionalString{
				Val: *update.Name,
			}
		}

		if update.ReplicaN != nil {
			pUpdate.ReplicaN = &internal.OptionalInt64{
				Val: int64(*update.ReplicaN),
			}
		}

		if update.Duration != nil {
			pUpdate.Duration = &internal.OptionalInt64{
				Val: int64(*update.Duration),
			}
		}

		if update.ShardGroupDuration != nil {
			pUpdate.ShardGroupDuration = &internal.OptionalInt64{
				Val: int64(*update.ShardGroupDuration),
			}
		}
	}

	return pUpdate
}

func nano2time(nano int64) time.Time {
	return time.Unix(nano/nanoi64, nano%nanoi64)
}
