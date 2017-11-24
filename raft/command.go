package raft

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
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
	fsmCmdTypeMetaDropUser                          = 8
	fsmCmdTypeMetaUpdateUser                        = 9
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
	fsmCmdTypeMetaDatabase               = 102
	fsmCmdTypeMetaDatabases              = 103
	fsmCmdTypeMetaRetentionPolicy        = 104
	fsmCmdTypeMetaUsers                  = 105
	fsmCmdTypeMetaUser                   = 106
	fsmCmdTypeMetaUserPrivileges         = 107
	fsmCmdTypeMetaUserPrivilege          = 108
	fsmCmdTypeMetaAdminUserExists        = 109
	fsmCmdTypeMetaAuthenticate           = 110
	fsmCmdTypeMetaShardGroupsByTimeRange = 111

	fsmCmdTypeDataCreateShard           = 51
	fsmCmdTypeDataWriteToShard          = 52
	fsmCmdTypeDataDeleteShard           = 53
	fsmCmdTypeDataDeleteDatabase        = 54
	fsmCmdTypeDataDeleteRetentionPolicy = 55
	fsmCmdTypeDataDeleteMeasurement     = 56
	fsmCmdTypeDataDeleteSeries          = 57

	fsmCmdTypeDataShardGroup              = 151
	fsmCmdTypeDataShardIDs                = 152
	fsmCmdTypeDataMeasurementNames        = 153
	fsmCmdTypeDataMeasurementsCardinality = 154
	fsmCmdTypeDataSeriesCardinality       = 155
	fsmCmdTypeDataTagKeys                 = 156
	fsmCmdTypeDataTagValues               = 157
)

func (f fsmCmdType) isLocalQuery() bool {
	return f/100 > 0
}

func applyCmd(r *Raft, cmdType fsmCmdType, cmd proto.Message) (interface{}, error) {
	b, id, err := fsmCmdMarshal(cmdType, cmd)
	if err != nil {
		return nil, err
	}

	if id != nil {
		r.queryMgr.Add(*id)
	}

	// TODO configurable
	f := r.r.Apply(b, 10*time.Second)

	if err := f.Error(); err != nil {
		return nil, err
	}

	resp := (f.Response()).(fsmCmdResponse)
	return resp.res, resp.err
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

func nano2time(nano int64) time.Time {
	return time.Unix(nano/nanoi64, nano%nanoi64)
}
