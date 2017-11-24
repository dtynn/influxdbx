package raft

import (
	"time"

	"github.com/dtynn/influxdbx/coordinator"
	"github.com/dtynn/influxdbx/raft/internal"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

var (
	emptyMetaUser = &meta.UserInfo{}

	_ coordinator.MetaClient = (*wrappedMetaClient)(nil)
)

type wrappedMetaClient struct {
	r *Raft
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
	cmd := &internal.MetaAcquireLeaseCmd{
		Name: name,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaAcquireLease, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.Lease), nil
}

// ClusterID returns the ID of the cluster it's connected to.
func (w *wrappedMetaClient) ClusterID() uint64 {
	return w.r.MetaClient.ClusterID()
}

func (w *wrappedMetaClient) Database(name string) (*meta.DatabaseInfo, error) {
	cmd := &internal.MetaDatabaseCmd{
		Name: name,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaDatabase, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) Databases() ([]meta.DatabaseInfo, error) {
	res, err := applyCmd(w.r, fsmCmdTypeMetaDatabases, nil)
	if err != nil {
		return nil, err
	}

	return res.([]meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) CreateDatabase(name string) (*meta.DatabaseInfo, error) {
	cmd := &internal.MetaCreateDatabaseCmd{
		Name: name,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaCreateDatabase, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) CreateDatabaseWithRetentionPolicy(name string, spec *meta.RetentionPolicySpec) (*meta.DatabaseInfo, error) {
	cmd := &internal.MetaCreateDatabaseWithRetentionPolicyCmd{
		Name: name,
		Spec: retentionPolicySpecMeta2Proto(spec),
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaCreateDatabaseWithRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.DatabaseInfo), nil
}

func (w *wrappedMetaClient) DropDatabase(name string) error {
	cmd := &internal.MetaDropDatabaseCmd{
		Name: name,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropDatabase, cmd)

	return err
}

func (w *wrappedMetaClient) CreateRetentionPolicy(database string, spec *meta.RetentionPolicySpec, makeDefault bool) (*meta.RetentionPolicyInfo, error) {
	cmd := &internal.MetaCreateRetentionPolicyCmd{
		Database:    database,
		Spec:        retentionPolicySpecMeta2Proto(spec),
		MakeDefault: makeDefault,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaCreateRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.RetentionPolicyInfo), nil
}

func (w *wrappedMetaClient) RetentionPolicy(database, name string) (*meta.RetentionPolicyInfo, error) {
	cmd := &internal.MetaRetentionPolicyCmd{
		Database: database,
		Name:     name,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaRetentionPolicy, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.RetentionPolicyInfo), nil
}

func (w *wrappedMetaClient) DropRetentionPolicy(database, name string) error {
	cmd := &internal.MetaDropRetentionPolicyCmd{
		Database: database,
		Name:     name,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropRetentionPolicy, cmd)

	return err
}

func (w *wrappedMetaClient) UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error {
	cmd := &internal.MetaUpdateRetentionPolicyCmd{
		Database:    database,
		Name:        name,
		Update:      retentionPolicyUpdateMeta2Proto(rpu),
		MakeDefault: makeDefault,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaUpdateRetentionPolicy, cmd)

	return err
}

func (w *wrappedMetaClient) Users() ([]meta.UserInfo, error) {
	res, err := applyCmd(w.r, fsmCmdTypeMetaUsers, nil)
	if err != nil {
		return nil, err
	}

	return res.([]meta.UserInfo), nil
}

func (w *wrappedMetaClient) User(name string) (meta.User, error) {
	cmd := &internal.MetaUserCmd{
		Name: name,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaUser, cmd)
	if err != nil {
		return nil, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) CreateUser(name, password string, admin bool) (meta.User, error) {
	cmd := &internal.MetaCreateUserCmd{
		Name:     name,
		Password: password,
		Admin:    admin,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaCreateUser, cmd)
	if err != nil {
		return emptyMetaUser, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) UpdateUser(name, password string) error {
	cmd := &internal.MetaUpdateUserCmd{
		Name:     name,
		Password: password,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaUpdateUser, cmd)

	return err
}

func (w *wrappedMetaClient) DropUser(name string) error {
	cmd := &internal.MetaDropUserCmd{
		Name: name,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropUser, cmd)

	return err
}

func (w *wrappedMetaClient) SetPrivilege(username, database string, p influxql.Privilege) error {
	cmd := &internal.MetaSetPrivilegeCmd{
		Username:  username,
		Database:  database,
		Privilege: int64(p),
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaSetPrivilege, cmd)

	return err
}

func (w *wrappedMetaClient) SetAdminPrivilege(username string, admin bool) error {
	cmd := &internal.MetaSetAdminPrivilegeCmd{
		Username: username,
		Admin:    admin,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaSetAdminPrivilege, cmd)

	return err
}

func (w *wrappedMetaClient) UserPrivileges(username string) (map[string]influxql.Privilege, error) {
	cmd := &internal.MetaUserPrivilegesCmd{
		Username: username,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaUserPrivileges, cmd)
	if err != nil {
		return nil, err
	}

	return res.(map[string]influxql.Privilege), nil
}

func (w *wrappedMetaClient) UserPrivilege(username, database string) (*influxql.Privilege, error) {
	cmd := &internal.MetaUserPrivilegeCmd{
		Username: username,
		Database: database,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaUserPrivilege, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*influxql.Privilege), nil
}

func (w *wrappedMetaClient) AdminUserExists() (bool, error) {
	res, err := applyCmd(w.r, fsmCmdTypeMetaAdminUserExists, nil)
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (w *wrappedMetaClient) Authenticate(username, password string) (meta.User, error) {
	cmd := &internal.MetaAuthenticateCmd{
		Username: username,
		Password: password,
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaAuthenticate, cmd)
	if err != nil {
		return nil, err
	}

	return res.(meta.User), nil
}

func (w *wrappedMetaClient) ShardGroupsByTimeRange(database, policy string, min, max time.Time) ([]meta.ShardGroupInfo, error) {
	cmd := &internal.MetaShardGroupsByTimeRangeCmd{
		Database: database,
		Policy:   policy,
		Tmin:     min.UnixNano(),
		Tmax:     max.UnixNano(),
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaShardGroupsByTimeRange, cmd)
	if err != nil {
		return nil, err
	}

	return res.([]meta.ShardGroupInfo), nil
}

func (w *wrappedMetaClient) DropShard(id uint64) error {
	cmd := &internal.MetaDropShardCmd{
		Id: id,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropShard, cmd)

	return err
}

func (w *wrappedMetaClient) PruneShardGroups() error {
	_, err := applyCmd(w.r, fsmCmdTypeMetaPruneShardGroups, nil)

	return err
}

func (w *wrappedMetaClient) CreateShardGroup(database, policy string, timestamp time.Time) (*meta.ShardGroupInfo, error) {
	cmd := &internal.MetaCreateShardGroupCmd{
		Database:  database,
		Policy:    policy,
		Timestamp: timestamp.UnixNano(),
	}

	res, err := applyCmd(w.r, fsmCmdTypeMetaCreateShardGroup, cmd)
	if err != nil {
		return nil, err
	}

	return res.(*meta.ShardGroupInfo), nil
}

func (w *wrappedMetaClient) DeleteShardGroup(database, policy string, id uint64) error {
	cmd := &internal.MetaDeleteShardGroupCmd{
		Database: database,
		Policy:   policy,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDeleteShardGroup, cmd)

	return err
}

func (w *wrappedMetaClient) PrecreateShardGroups(from, to time.Time) error {
	cmd := &internal.MetaPrecreateShardGroupsCmd{
		From: from.UnixNano(),
		To:   to.UnixNano(),
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaPrecreateShardGroups, cmd)

	return err
}

func (w *wrappedMetaClient) CreateContinuousQuery(database, name, query string) error {
	cmd := &internal.MetaCreateContinuousQueryCmd{
		Database: database,
		Name:     name,
		Query:    query,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaCreateContinuousQuery, cmd)
	return err
}

func (w *wrappedMetaClient) DropContinuousQuery(database, name string) error {
	cmd := &internal.MetaDropContinuousQueryCmd{
		Database: database,
		Name:     name,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropContinuousQuery, cmd)

	return err
}

func (w *wrappedMetaClient) CreateSubscription(database, rp, name, mode string, destinations []string) error {
	cmd := &internal.MetaCreateSubscriptionCmd{
		Database:     database,
		Rp:           rp,
		Name:         name,
		Mode:         mode,
		Destinations: destinations,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaCreateSubscription, cmd)

	return err
}

func (w *wrappedMetaClient) DropSubscription(database, rp, name string) error {
	cmd := &internal.MetaDropSubscriptionCmd{
		Database: database,
		Rp:       rp,
		Name:     name,
	}

	_, err := applyCmd(w.r, fsmCmdTypeMetaDropSubscription, cmd)

	return err
}

func (w *wrappedMetaClient) WaitForDataChanged() chan struct{} {
	return w.r.MetaClient.WaitForDataChanged()
}

func (w *wrappedMetaClient) MarshalBinary() ([]byte, error) {
	return w.r.MetaClient.MarshalBinary()
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

func userInfoMeta2Proto(ui *meta.UserInfo) *internal.UserInfo {
	if ui == nil {
		return nil
	}

	res := &internal.UserInfo{
		Name:       ui.Name,
		Hash:       ui.Hash,
		Admin:      ui.Admin,
		Privileges: make(map[string]int64, len(ui.Privileges)),
	}

	for db, p := range ui.Privileges {
		res.Privileges[db] = int64(p)
	}

	return res
}

func userInfoProto2Meta(ui *internal.UserInfo) *meta.UserInfo {
	if ui == nil {
		return nil
	}

	res := &meta.UserInfo{
		Name:       ui.GetName(),
		Hash:       ui.GetHash(),
		Admin:      ui.GetAdmin(),
		Privileges: make(map[string]influxql.Privilege, len(ui.Privileges)),
	}

	for db, p := range ui.GetPrivileges() {
		res.Privileges[db] = influxql.Privilege(p)
	}

	return res
}
