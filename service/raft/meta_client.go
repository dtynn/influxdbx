package raft

import (
	"time"

	"github.com/dtynn/influxdbx/service/raft/internal"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

var (
	emptyMetaUser = &meta.UserInfo{}
)

type WrappedMetaClient struct {
	r        *Raft
	queryMgr *localQueryManager
}

func NewWrappedMetaClient(r *Raft) *WrappedMetaClient {
	return &WrappedMetaClient{
		r: r,
	}
}

func (w *WrappedMetaClient) applyCmd(cmdType fsmCmdType, cmd proto.Message) (interface{}, error) {
	b, id, err := fsmCmdMarshal(cmdType, cmd)
	if err != nil {
		return nil, err
	}

	if id != nil {
		w.queryMgr.Add(*id)
	}

	// TODO configurable
	f := w.r.r.Apply(b, 10*time.Second)

	if err := f.Error(); err != nil {
		return nil, err
	}

	resp := (f.Response()).(fsmCmdResponse)
	return resp.res, resp.err
}

// AcquireLease
// A lease is a logical concept that can be used by anything that needs to limit
// execution to a single node.  E.g., the CQ service on all nodes may ask for
// the "ContinuousQuery" lease. Only the node that acquires it will run CQs.
// NOTE: Leases are not managed through the CP system and are not fully
// consistent.  Any actions taken after acquiring a lease must be idempotent.
// TODO read code
// services/continuous_querier Service.backgroundLoop
func (w *WrappedMetaClient) AcquireLease(name string) (*meta.Lease, error) {
	cmd := &internal.AcquireLeaseCmd{
		Name: name,
	}

	resp, err := w.applyCmd(fsmCmdTypeAcquireLease, cmd)
	if err != nil {
		return nil, err
	}

	return resp.(*meta.Lease), nil
}

// ClusterID returns the ID of the cluster it's connected to.
func (w *WrappedMetaClient) ClusterID() uint64 {
	return w.r.MetaClient.ClusterID()
}

func (w *WrappedMetaClient) Database(name string) (*meta.DatabaseInfo, error) {
	return nil, nil
}

func (w *WrappedMetaClient) Databases() []meta.DatabaseInfo {
	return nil
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

func (w *WrappedMetaClient) DropDatabase(name string) error {
	return nil
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

func (w *WrappedMetaClient) RetentionPolicy(database, name string) (rpi *meta.RetentionPolicyInfo, err error) {
	return nil, nil
}

func (w *WrappedMetaClient) DropRetentionPolicy(database, name string) error {
	return nil
}

func (w *WrappedMetaClient) UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error {
	return nil
}

func (w *WrappedMetaClient) Users() []meta.UserInfo {
	return nil
}

func (w *WrappedMetaClient) User(name string) (meta.User, error) {
	return nil, nil
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

func (w *WrappedMetaClient) UpdateUser(name, password string) error {
	return nil
}

func (w *WrappedMetaClient) DropUser(name string) error {
	return nil
}

func (w *WrappedMetaClient) SetPrivilege(username, database string, p influxql.Privilege) error {
	return nil
}

func (w *WrappedMetaClient) SetAdminPrivilege(username string, admin bool) error {
	return nil
}

func (w *WrappedMetaClient) UserPrivileges(username string) (map[string]influxql.Privilege, error) {
	return nil, nil
}

func (w *WrappedMetaClient) UserPrivilege(username, database string) (*influxql.Privilege, error) {
	return nil, nil
}

func (w *WrappedMetaClient) AdminUserExists() bool {
	return false
}

func (w *WrappedMetaClient) Authenticate(username, password string) (meta.User, error) {
	return nil, nil
}

func (w *WrappedMetaClient) UserCount() int {
	return 0
}

func (w *WrappedMetaClient) ShardIDs() []uint64 {
	return nil
}

func (w *WrappedMetaClient) ShardGroupsByTimeRange(database, policy string, min, max time.Time) (a []meta.ShardGroupInfo, err error) {
	return nil, nil
}

func (w *WrappedMetaClient) ShardsByTimeRange(sources influxql.Sources, tmin, tmax time.Time) (a []meta.ShardInfo, err error) {
	return nil, nil
}

func (w *WrappedMetaClient) DropShard(id uint64) error {
	return nil
}

func (w *WrappedMetaClient) PruneShardGroups() error {
	return nil
}

func (w *WrappedMetaClient) CreateShardGroup(database, policy string, timestamp time.Time) (*meta.ShardGroupInfo, error) {
	return nil, nil
}

func (w *WrappedMetaClient) DeleteShardGroup(database, policy string, id uint64) error {
	return nil
}

func (w *WrappedMetaClient) PrecreateShardGroups(from, to time.Time) error {
	return nil
}

func (w *WrappedMetaClient) ShardOwner(shardID uint64) (database, policy string, sgi *meta.ShardGroupInfo) {
	return "", "", nil
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

func (w *WrappedMetaClient) DropContinuousQuery(database, name string) error {
	return nil
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

func (w *WrappedMetaClient) DropSubscription(database, rp, name string) error {
	return nil
}

func (w *WrappedMetaClient) SetData(data *meta.Data) error {
	return nil
}

func (w *WrappedMetaClient) Data() meta.Data {
	return w.r.MetaClient.Data()
}

func (w *WrappedMetaClient) WaitForDataChanged() chan struct{} {
	return w.r.MetaClient.WaitForDataChanged()
}

func (w *WrappedMetaClient) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (w *WrappedMetaClient) Load() error {
	return w.r.MetaClient.Load()
}
