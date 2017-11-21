package raft

import (
	"fmt"
	"time"

	"github.com/dtynn/influxdbx/service/raft/internal"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb/services/meta"
)

type fsmCmdType byte

const (
	fsmCmdTypeUnknown fsmCmdType = 0

	fsmCmdTypeCreateContinuousQuery = 11
	fsmCmdTypeDropContinuousQuery   = 12

	fsmCmdTypeCreateDatabase                    = 21
	fsmCmdTypeCreateDatabaseWithRetentionPolicy = 22
	fsmCmdTypeDropDatabase                      = 23

	fsmCmdTypeCreateRetentionPolicy = 31
	fsmCmdTypeDropRetentionPolicy   = 32
	fsmCmdTypeUpdateRetentionPolicy = 33

	fsmCmdTypeCreateSubscription = 41
	fsmCmdTypeDropSubscription   = 42

	fsmCmdTypeCreateUser        = 51
	fsmCmdTypeDropUser          = 52
	fsmCmdTypeSetAdminPrivilege = 53
	fsmCmdTypeSetPrivilege      = 54
	fsmCmdTypeUpdateUser        = 55

	fsmCmdTypeDropShard = 61
)

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
