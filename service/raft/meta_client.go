package raft

import (
	"time"

	"github.com/dtynn/influxdbx/service/raft/internal"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb/services/meta"
)

var (
	emptyMetaUser = &meta.UserInfo{}
)

type WrappedMetaClient struct {
	r *Raft
}

func NewWrappedMetaClient(r *Raft) *WrappedMetaClient {
	return &WrappedMetaClient{
		r: r,
	}
}

func (w *WrappedMetaClient) applyCmd(fsmCmdType byte, cmd proto.Message) (interface{}, error) {
	b, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	bs := make([]byte, len(b)+1)
	bs[0] = fsmCmdType
	copy(bs[1:], b)

	// TODO configurable
	f := w.r.r.Apply(bs, 10*time.Second)

	if err := f.Error(); err != nil {
		return nil, err
	}

	resp := (f.Response()).(fsmCmdResponse)
	return resp.res, resp.err
}

func (w *WrappedMetaClient) CreateContinuousQuery(database, name, query string) error {
	cmd := &internal.CreateContinuousQueryCmd{
		Database: database,
		Name:     name,
		Query:    query,
	}

	_, err := w.applyCmd(fsmCmdTypeCreateContinuousQuery, cmd)
	return err
}

func (w *WrappedMetaClient) CreateDatabase(name string) (*meta.DatabaseInfo, error) {
	cmd := &internal.CreateDatabaseCmd{
		Name: name,
	}

	res, err := w.applyCmd(fsmCmdTypeCreateDatabase, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *WrappedMetaClient) CreateDatabaseWithRetentionPolicy(name string, spec *meta.RetentionPolicySpec) (*meta.DatabaseInfo, error) {
	cmd := &internal.CreateDatabaseWithRetentionPolicyCmd{
		Name: name,
		Spec: retentionPolicySpecMeta2Proto(spec),
	}

	res, err := w.applyCmd(fsmCmdTypeCreateDatabaseWithRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *WrappedMetaClient) CreateRetentionPolicy(database string, spec *meta.RetentionPolicySpec, makeDefault bool) (*meta.RetentionPolicyInfo, error) {
	cmd := &internal.CreateRetentionPolicyCmd{
		Database:    database,
		Spec:        retentionPolicySpecMeta2Proto(spec),
		MakeDefault: makeDefault,
	}

	res, err := w.applyCmd(fsmCmdTypeCreateRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.RetentionPolicyInfo), nil
}

func (w *WrappedMetaClient) CreateSubscription(database, rp, name, mode string, destinations []string) error {
	cmd := &internal.CreateSubscriptionCmd{
		Database:     database,
		Rp:           rp,
		Name:         name,
		Mode:         mode,
		Destinations: destinations,
	}

	_, err := w.applyCmd(fsmCmdTypeCreateSubscription, cmd)
	return err
}

func (w *WrappedMetaClient) CreateUser(name, password string, admin bool) (meta.User, error) {
	cmd := &internal.CreateUserCmd{
		Name:     name,
		Password: password,
		Admin:    admin,
	}

	res, err := w.applyCmd(fsmCmdTypeCreateUser, cmd)
	if err != nil {
		return emptyMetaUser, err
	}

	return res.(meta.User), nil
}

// func (w *WrappedMetaClient) Database(name string) *meta.DatabaseInfo {

// }

// func (w *WrappedMetaClient) Databases() []meta.DatabaseInfo {

// }

// func (w *WrappedMetaClient) DropShard(id uint64) error {

// }

// func (w *WrappedMetaClient) DropContinuousQuery(database, name string) error {

// }

// func (w *WrappedMetaClient) DropDatabase(name string) error {

// }

// func (w *WrappedMetaClient) DropRetentionPolicy(database, name string) error {

// }

// func (w *WrappedMetaClient) DropSubscription(database, rp, name string) error {

// }

// func (w *WrappedMetaClient) DropUser(name string) error {

// }

// func (w *WrappedMetaClient) RetentionPolicy(database, name string) (rpi *meta.RetentionPolicyInfo, err error) {

// }

// func (w *WrappedMetaClient) SetAdminPrivilege(username string, admin bool) error {

// }

// func (w *WrappedMetaClient) SetPrivilege(username, database string, p influxql.Privilege) error {

// }

// func (w *WrappedMetaClient) ShardGroupsByTimeRange(database, policy string, min, max time.Time) (a []meta.ShardGroupInfo, err error) {

// }

// func (w *WrappedMetaClient) UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error {

// }

// func (w *WrappedMetaClient) UpdateUser(name, password string) error {

// }

// func (w *WrappedMetaClient) UserPrivilege(username, database string) (*influxql.Privilege, error) {

// }

// func (w *WrappedMetaClient) UserPrivileges(username string) (map[string]influxql.Privilege, error) {

// }

// func (w *WrappedMetaClient) Users() []meta.UserInfo {

// }
