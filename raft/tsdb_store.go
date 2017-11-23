package raft

import (
	"io"
	"time"
)

type wrappedTSDBStore struct {
	r *Raft
}

func (w *wrappedTSDBStore) BackupShard(id uint64, since time.Time, w io.Writer) error {

}

func (w *wrappedTSDBStore) Close() error {

}

func (w *wrappedTSDBStore) CreateShard(database, policy string, shardID uint64, enabled bool) error {

}

func (w *wrappedTSDBStore) CreateShardSnapshot(id uint64) (string, error) {

}

func (w *wrappedTSDBStore) Databases() []string {

}

func (w *wrappedTSDBStore) DeleteDatabase(name string) error {

}

func (w *wrappedTSDBStore) DeleteMeasurement(database, name string) error {

}

func (w *wrappedTSDBStore) DeleteRetentionPolicy(database, name string) error {

}

func (w *wrappedTSDBStore) DeleteSeries(database string, sources []influxql.Source, condition influxql.Expr) error {

}

func (w *wrappedTSDBStore) DeleteShard(id uint64) error {

}

func (w *wrappedTSDBStore) DiskSize() (int64, error) {

}

func (w *wrappedTSDBStore) ExpandSources(sources influxql.Sources) (influxql.Sources, error) {

}

func (w *wrappedTSDBStore) ImportShard(id uint64, r io.Reader) error {

}

func (w *wrappedTSDBStore) MeasurementSeriesCounts(database string) (measuments int, series int) {

}

func (w *wrappedTSDBStore) MeasurementsCardinality(database string) (int64, error) {

}

func (w *wrappedTSDBStore) MeasurementNames(auth query.Authorizer, database string, cond influxql.Expr) ([][]byte, error) {

}

func (w *wrappedTSDBStore) Open() error {

}

func (w *wrappedTSDBStore) Path() string {

}

func (w *wrappedTSDBStore) RestoreShard(id uint64, r io.Reader) error {

}

func (w *wrappedTSDBStore) SeriesCardinality(database string) (int64, error) {

}

func (w *wrappedTSDBStore) SetShardEnabled(shardID uint64, enabled bool) error {

}

func (w *wrappedTSDBStore) Shard(id uint64) *tsdb.Shard {

}

func (w *wrappedTSDBStore) ShardGroup(ids []uint64) tsdb.ShardGroup {

}

func (w *wrappedTSDBStore) ShardIDs() []uint64 {

}

func (w *wrappedTSDBStore) ShardN() int {

}

func (w *wrappedTSDBStore) ShardRelativePath(id uint64) (string, error) {

}

func (w *wrappedTSDBStore) Shards(ids []uint64) []*tsdb.Shard {

}

func (w *wrappedTSDBStore) Statistics(tags map[string]string) []models.Statistic {

}

func (w *wrappedTSDBStore) TagKeys(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagKeys, error) {

}

func (w *wrappedTSDBStore) TagValues(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagValues, error) {

}

func (w *wrappedTSDBStore) WithLogger(log *zap.Logger) {

}

func (w *wrappedTSDBStore) WriteToShard(shardID uint64, points []models.Point) error {

}
