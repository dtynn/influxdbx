package raft

import (
	"fmt"
	"time"

	"github.com/dtynn/influxdbx/raft/internal"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

var (
	emptyMetaUser = &meta.UserInfo{}
)

type wrappedMetaClient struct {
	r *Raft
}

func (w *wrappedMetaClient) applyCmd(cmdType fsmCmdType, cmd proto.Message) (interface{}, error) {
	b, id, err := fsmCmdMarshal(cmdType, cmd)
	if err != nil {
		return nil, err
	}

	if id != nil {
		w.r.queryMgr.Add(*id)
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
func (w *wrappedMetaClient) AcquireLease(name string) (*meta.Lease, error) {
	cmd := &internal.AcquireLeaseCmd{
		Name: name,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaAcquireLease, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.Lease), nil
}

// ClusterID returns the ID of the cluster it's connected to.
func (w *wrappedMetaClient) ClusterID() uint64 {
	return w.r.MetaClient.ClusterID()
}

func (w *wrappedMetaClient) Database(name string) *meta.DatabaseInfo {
	cmd := &internal.DatabaseCmd{
		Name: name,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaDatabase, cmd)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.Database: %s", err))
		return nil
	}

	return res.(*meta.DatabaseInfo)
}

func (w *wrappedMetaClient) Databases() []meta.DatabaseInfo {
	res, err := w.applyCmd(fsmCmdTypeMetaDatabases, nil)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.Databases: %s", err))
		return nil
	}

	return res.([]meta.DatabaseInfo)
}

func (w *wrappedMetaClient) CreateDatabase(name string) (*meta.DatabaseInfo, error) {
	cmd := &internal.CreateDatabaseCmd{
		Name: name,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaCreateDatabase, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) CreateDatabaseWithRetentionPolicy(name string, spec *meta.RetentionPolicySpec) (*meta.DatabaseInfo, error) {
	cmd := &internal.CreateDatabaseWithRetentionPolicyCmd{
		Name: name,
		Spec: retentionPolicySpecMeta2Proto(spec),
	}

	res, err := w.applyCmd(fsmCmdTypeMetaCreateDatabaseWithRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) DropDatabase(name string) error {
	cmd := &internal.DropDatabaseCmd{
		Name: name,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropDatabase, cmd)

	return err
}

func (w *wrappedMetaClient) CreateRetentionPolicy(database string, spec *meta.RetentionPolicySpec, makeDefault bool) (*meta.RetentionPolicyInfo, error) {
	cmd := &internal.CreateRetentionPolicyCmd{
		Database:    database,
		Spec:        retentionPolicySpecMeta2Proto(spec),
		MakeDefault: makeDefault,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaCreateRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.RetentionPolicyInfo), nil
}

func (w *wrappedMetaClient) RetentionPolicy(database, name string) (*meta.RetentionPolicyInfo, error) {
	cmd := &internal.RetentionPolicyCmd{
		Database: database,
		Name:     name,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.RetentionPolicyInfo), nil
}

func (w *wrappedMetaClient) DropRetentionPolicy(database, name string) error {
	cmd := &internal.DropRetentionPolicyCmd{
		Database: database,
		Name:     name,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropRetentionPolicy, cmd)

	return err
}

func (w *wrappedMetaClient) UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error {
	cmd := &internal.UpdateRetentionPolicyCmd{
		Database:    database,
		Name:        name,
		Update:      retentionPolicyUpdateMeta2Proto(rpu),
		MakeDefault: makeDefault,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaUpdateRetentionPolicy, cmd)

	return err
}

func (w *wrappedMetaClient) Users() []meta.UserInfo {
	res, err := w.applyCmd(fsmCmdTypeMetaUsers, nil)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.Users: %s", err))
		return nil
	}

	return res.([]meta.UserInfo)
}

func (w *wrappedMetaClient) User(name string) (meta.User, error) {
	cmd := &internal.UserCmd{
		Name: name,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaUser, cmd)
	if err != nil {
		return nil, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) CreateUser(name, password string, admin bool) (meta.User, error) {
	cmd := &internal.CreateUserCmd{
		Name:     name,
		Password: password,
		Admin:    admin,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaCreateUser, cmd)
	if err != nil {
		return emptyMetaUser, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) UpdateUser(name, password string) error {
	cmd := &internal.UpdateUserCmd{
		Name:     name,
		Password: password,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaUpdateUser, cmd)

	return err
}

func (w *wrappedMetaClient) DropUser(name string) error {
	cmd := &internal.DropUserCmd{
		Name: name,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropUser, cmd)

	return err
}

func (w *wrappedMetaClient) SetPrivilege(username, database string, p influxql.Privilege) error {
	cmd := &internal.SetPrivilegeCmd{
		Username:  username,
		Database:  database,
		Privilege: int64(p),
	}

	_, err := w.applyCmd(fsmCmdTypeMetaSetPrivilege, cmd)

	return err
}

func (w *wrappedMetaClient) SetAdminPrivilege(username string, admin bool) error {
	cmd := &internal.SetAdminPrivilegeCmd{
		Username: username,
		Admin:    admin,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaSetAdminPrivilege, cmd)

	return err
}

func (w *wrappedMetaClient) UserPrivileges(username string) (map[string]influxql.Privilege, error) {
	cmd := &internal.UserPrivilegesCmd{
		Username: username,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaUserPrivileges, cmd)
	if err != nil {
		return nil, err
	}

	return res.(map[string]influxql.Privilege), nil
}

func (w *wrappedMetaClient) UserPrivilege(username, database string) (*influxql.Privilege, error) {
	cmd := &internal.UserPrivilegeCmd{
		Username: username,
		Database: database,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaUserPrivilege, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*influxql.Privilege), nil
}

func (w *wrappedMetaClient) AdminUserExists() bool {
	res, err := w.applyCmd(fsmCmdTypeMetaAdminUserExists, nil)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.AdminUserExists: %s", err))
		return false
	}

	return res.(bool)
}

func (w *wrappedMetaClient) Authenticate(username, password string) (meta.User, error) {
	cmd := &internal.AuthenticateCmd{
		Username: username,
		Password: password,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaAuthenticate, cmd)
	if err != nil {
		return nil, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) UserCount() int {
	res, err := w.applyCmd(fsmCmdTypeMetaUserCount, nil)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.UserCount: %s", err))
		return 0
	}

	return res.(int)
}

func (w *wrappedMetaClient) ShardIDs() []uint64 {
	res, err := w.applyCmd(fsmCmdTypeMetaShardIDs, nil)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.ShardIDs: %s", err))
		return nil
	}

	return res.([]uint64)
}

func (w *wrappedMetaClient) ShardGroupsByTimeRange(database, policy string, min, max time.Time) ([]meta.ShardGroupInfo, error) {
	cmd := &internal.ShardGroupsByTimeRangeCmd{
		Database: database,
		Policy:   policy,
		Tmin:     min.UnixNano(),
		Tmax:     max.UnixNano(),
	}

	res, err := w.applyCmd(fsmCmdTypeMetaShardGroupsByTimeRange, cmd)
	if err != nil {
		return nil, err
	}

	return res.([]meta.ShardGroupInfo), nil
}

func (w *wrappedMetaClient) ShardsByTimeRange(sources influxql.Sources, tmin, tmax time.Time) ([]meta.ShardInfo, error) {
	sourcesData, err := sources.MarshalBinary()
	if err != nil {
		return nil, err
	}

	cmd := &internal.ShardsByTimeRangeCmd{
		SourceData: sourcesData,
		Tmin:       tmin.UnixNano(),
		Tmax:       tmax.UnixNano(),
	}

	res, err := w.applyCmd(fsmCmdTypeMetaShardsByTimeRange, cmd)
	if err != nil {
		return nil, err
	}

	return res.([]meta.ShardInfo), nil
}

func (w *wrappedMetaClient) DropShard(id uint64) error {
	cmd := &internal.DropShardCmd{
		Id: id,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropShard, cmd)

	return err
}

func (w *wrappedMetaClient) PruneShardGroups() error {
	_, err := w.applyCmd(fsmCmdTypeMetaPruneShardGroups, nil)

	return err
}

func (w *wrappedMetaClient) CreateShardGroup(database, policy string, timestamp time.Time) (*meta.ShardGroupInfo, error) {
	cmd := &internal.CreateShardGroupCmd{
		Database:  database,
		Policy:    policy,
		Timestamp: timestamp.UnixNano(),
	}

	res, err := w.applyCmd(fsmCmdTypeMetaCreateShardGroup, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.ShardGroupInfo), nil
}

func (w *wrappedMetaClient) DeleteShardGroup(database, policy string, id uint64) error {
	cmd := &internal.DeleteShardGroupCmd{
		Database: database,
		Policy:   policy,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDeleteShardGroup, cmd)

	return err
}

func (w *wrappedMetaClient) PrecreateShardGroups(from, to time.Time) error {
	cmd := &internal.PrecreateShardGroupsCmd{
		From: from.UnixNano(),
		To:   to.UnixNano(),
	}

	_, err := w.applyCmd(fsmCmdTypeMetaPrecreateShardGroups, cmd)

	return err
}

type shardOwnerResult struct {
	database string
	policy   string
	info     *meta.ShardGroupInfo
}

func (w *wrappedMetaClient) ShardOwner(shardID uint64) (database, policy string, sgi *meta.ShardGroupInfo) {
	cmd := &internal.ShardOwnerCmd{
		Id: shardID,
	}

	res, err := w.applyCmd(fsmCmdTypeMetaShardOwner, cmd)
	if err != nil {
		w.r.logger.Warn(fmt.Sprintf("wrappedMetaClient.ShardOwner: %s", err))
		return "", "", nil
	}

	sor := res.(shardOwnerResult)

	return sor.database, sor.policy, sor.info
}

func (w *wrappedMetaClient) CreateContinuousQuery(database, name, query string) error {
	cmd := &internal.CreateContinuousQueryCmd{
		Database: database,
		Name:     name,
		Query:    query,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaCreateContinuousQuery, cmd)
	return err
}

func (w *wrappedMetaClient) DropContinuousQuery(database, name string) error {
	cmd := &internal.DropContinuousQueryCmd{
		Database: database,
		Name:     name,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropContinuousQuery, cmd)

	return err
}

func (w *wrappedMetaClient) CreateSubscription(database, rp, name, mode string, destinations []string) error {
	cmd := &internal.CreateSubscriptionCmd{
		Database:     database,
		Rp:           rp,
		Name:         name,
		Mode:         mode,
		Destinations: destinations,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaCreateSubscription, cmd)

	return err
}

func (w *wrappedMetaClient) DropSubscription(database, rp, name string) error {
	cmd := &internal.DropSubscriptionCmd{
		Database: database,
		Rp:       rp,
		Name:     name,
	}

	_, err := w.applyCmd(fsmCmdTypeMetaDropSubscription, cmd)

	return err
}

func (w *wrappedMetaClient) SetData(data *meta.Data) error {
	return w.r.MetaClient.SetData(data)
}

func (w *wrappedMetaClient) Data() meta.Data {
	return w.r.MetaClient.Data()
}

func (w *wrappedMetaClient) WaitForDataChanged() chan struct{} {
	return w.r.MetaClient.WaitForDataChanged()
}

func (w *wrappedMetaClient) MarshalBinary() ([]byte, error) {
	return w.r.MetaClient.MarshalBinary()
}

func (w *wrappedMetaClient) Load() error {
	return w.r.MetaClient.Load()
}
