package coordinator

import (
	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/query"
	"github.com/influxdata/influxdb/tsdb"
	"github.com/influxdata/influxql"
)

type TSDBStore interface {
	// coordinator.ShardGroup
	ShardGroup(ids []uint64) (tsdb.ShardGroup, error)

	// coordinator.StatementExecutor
	CreateShard(database, policy string, shardID uint64, enabled bool) error
	WriteToShard(shardID uint64, points []models.Point) error

	DeleteDatabase(name string) error
	DeleteMeasurement(database, name string) error
	DeleteRetentionPolicy(database, name string) error
	DeleteSeries(database string, sources []influxql.Source, condition influxql.Expr) error
	DeleteShard(id uint64) error

	MeasurementNames(auth query.Authorizer, database string, cond influxql.Expr) ([][]byte, error)
	TagKeys(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagKeys, error)
	TagValues(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagValues, error)

	SeriesCardinality(database string) (int64, error)
	MeasurementsCardinality(database string) (int64, error)

	// services/rentention
	ShardIDs() ([]uint64, error)
}
