package coordinator

import (
	"time"

	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

type MetaClient interface {
	// coordinator.MetaClient
	CreateContinuousQuery(database, name, query string) error
	CreateDatabase(name string) (*meta.DatabaseInfo, error)
	CreateDatabaseWithRetentionPolicy(name string, spec *meta.RetentionPolicySpec) (*meta.DatabaseInfo, error)
	CreateRetentionPolicy(database string, spec *meta.RetentionPolicySpec, makeDefault bool) (*meta.RetentionPolicyInfo, error)
	CreateSubscription(database, rp, name, mode string, destinations []string) error
	CreateUser(name, password string, admin bool) (meta.User, error)
	Database(name string) (*meta.DatabaseInfo, error)
	Databases() ([]meta.DatabaseInfo, error)
	DropShard(id uint64) error
	DropContinuousQuery(database, name string) error
	DropDatabase(name string) error
	DropRetentionPolicy(database, name string) error
	DropSubscription(database, rp, name string) error
	DropUser(name string) error
	RetentionPolicy(database, name string) (rpi *meta.RetentionPolicyInfo, err error)
	SetAdminPrivilege(username string, admin bool) error
	SetPrivilege(username, database string, p influxql.Privilege) error
	ShardGroupsByTimeRange(database, policy string, min, max time.Time) (a []meta.ShardGroupInfo, err error)
	UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error
	UpdateUser(name, password string) error
	UserPrivilege(username, database string) (*influxql.Privilege, error)
	UserPrivileges(username string) (map[string]influxql.Privilege, error)
	Users() ([]meta.UserInfo, error)

	// coordinator.PointsWriter
	CreateShardGroup(database, policy string, timestamp time.Time) (*meta.ShardGroupInfo, error)

	// services/continuous_querier
	AcquireLease(name string) (l *meta.Lease, err error)

	// services/httpd
	Authenticate(username, password string) (ui meta.User, err error)
	User(username string) (meta.User, error)
	AdminUserExists() (bool, error)

	// services/prcreator
	PrecreateShardGroups(now, cutoff time.Time) error

	// services/retention
	DeleteShardGroup(database, policy string, id uint64) error
	PruneShardGroups() error

	// services/snapshotter
	MarshalBinary() (data []byte, err error)

	// services/subscriber
	WaitForDataChanged() chan struct{}
}
