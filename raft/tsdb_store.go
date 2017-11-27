package raft

import (
	"github.com/dtynn/influxdbx/coordinator"
	"github.com/dtynn/influxdbx/raft/internal"
	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/query"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxdb/tsdb"
	"github.com/influxdata/influxql"
)

var (
	_ coordinator.TSDBStore = (*wrappedTSDBStore)(nil)
)

type wrappedTSDBStore struct {
	r *Raft
}

func (wts *wrappedTSDBStore) CreateShard(database, policy string, shardID uint64, enabled bool) error {
	cmd := &internal.DataCreateShardCmd{
		Database: database,
		Policy:   policy,
		ShardID:  shardID,
		Enabled:  enabled,
	}

	_, err := applyCmd(wts.r, fsmCmdTypeDataCreateShard, cmd)
	return err
}

func (wts *wrappedTSDBStore) WriteToShard(shardID uint64, points []models.Point) error {
	pointData, err := pointsData2Proto(points)
	if err != nil {
		return err
	}

	cmd := &internal.DataWriteToShardCmd{
		ShardID:   shardID,
		PointData: pointData,
	}

	_, err = applyCmd(wts.r, fsmCmdTypeDataWriteToShard, cmd)

	return err
}

func (wts *wrappedTSDBStore) DeleteShard(id uint64) error {
	cmd := &internal.DataDeleteShardCmd{
		Id: id,
	}

	_, err := applyCmd(wts.r, fsmCmdTypeDataDeleteShard, cmd)

	return err
}

func (wts *wrappedTSDBStore) ShardGroup(ids []uint64) (tsdb.ShardGroup, error) {
	cmd := &internal.DataShardGroupCmd{
		Ids: ids,
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataShardGroup, cmd)
	if err != nil {
		return nil, err
	}

	return res.(tsdb.ShardGroup), nil
}

func (wts *wrappedTSDBStore) ShardIDs() ([]uint64, error) {
	res, err := applyCmd(wts.r, fsmCmdTypeDataShardIDs, nil)
	if err != nil {
		return nil, err
	}

	return res.([]uint64), nil
}

func (wts *wrappedTSDBStore) DeleteDatabase(name string) error {
	cmd := &internal.DataDeleteDatabaseCmd{
		Name: name,
	}

	_, err := applyCmd(wts.r, fsmCmdTypeDataDeleteDatabase, cmd)

	return err
}

func (wts *wrappedTSDBStore) DeleteRetentionPolicy(database, name string) error {
	cmd := &internal.DataDeleteRetentionPolicyCmd{
		Database: database,
		Name:     name,
	}

	_, err := applyCmd(wts.r, fsmCmdTypeDataDeleteRetentionPolicy, cmd)

	return err
}

func (wts *wrappedTSDBStore) DeleteMeasurement(database, name string) error {
	cmd := &internal.DataDeleteMeasurementCmd{
		Database: database,
		Name:     name,
	}

	_, err := applyCmd(wts.r, fsmCmdTypeDataDeleteMeasurement, cmd)

	return err
}

func (wts *wrappedTSDBStore) MeasurementNames(auth query.Authorizer, database string, cond influxql.Expr) ([][]byte, error) {
	cmd := &internal.DataMeasurementNamesCmd{
		User:      authorizer2ProtoUser(auth),
		Database:  database,
		Condition: cond.String(),
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataMeasurementNames, cmd)
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

func (wts *wrappedTSDBStore) MeasurementsCardinality(database string) (int64, error) {
	cmd := &internal.DataMeasurementsCardinalityCmd{
		Database: database,
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataMeasurementsCardinality, cmd)
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func (wts *wrappedTSDBStore) DeleteSeries(database string, sources []influxql.Source, condition influxql.Expr) error {
	sourcesData, err := influxql.Sources(sources).MarshalBinary()
	if err != nil {
		return err
	}

	cmd := &internal.DataDeleteSeriesCmd{
		Database:   database,
		SourceData: sourcesData,
	}

	_, err = applyCmd(wts.r, fsmCmdTypeDataDeleteSeries, cmd)

	return err
}

func (wts *wrappedTSDBStore) SeriesCardinality(database string) (int64, error) {
	cmd := &internal.DataSeriesCardinalityCmd{
		Database: database,
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataSeriesCardinality, cmd)
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func (wts *wrappedTSDBStore) TagKeys(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagKeys, error) {
	cmd := &internal.DataTagKeysCmd{
		User:      authorizer2ProtoUser(auth),
		ShardIDs:  shardIDs,
		Condition: cond.String(),
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataTagKeys, cmd)
	if err != nil {
		return nil, err
	}

	return res.([]tsdb.TagKeys), nil
}

func (wts *wrappedTSDBStore) TagValues(auth query.Authorizer, shardIDs []uint64, cond influxql.Expr) ([]tsdb.TagValues, error) {
	cmd := &internal.DataTagValuesCmd{
		User:      authorizer2ProtoUser(auth),
		ShardIDs:  shardIDs,
		Condition: cond.String(),
	}

	res, err := applyCmd(wts.r, fsmCmdTypeDataTagValues, cmd)
	if err != nil {
		return nil, err
	}

	return res.([]tsdb.TagValues), nil
}

func pointsData2Proto(points []models.Point) ([]byte, error) {
	size := 0
	for _, p := range points {
		size += p.StringSize() + 1
	}

	buf := make([]byte, 0, size)
	for i := range points {
		buf = points[i].AppendString(buf)
		buf = append(buf, '\n')
	}

	return buf, nil
}

func pointsProto2Data(buf []byte) ([]models.Point, error) {
	return models.ParsePoints(buf)
}

func authorizer2ProtoUser(auth query.Authorizer) *internal.UserInfo {
	if ui, ok := auth.(*meta.UserInfo); ok {
		return userInfoMeta2Proto(ui)
	}

	return nil
}

func protoUser2Authorizer(ui *internal.UserInfo) query.Authorizer {
	if ui == nil {
		return query.OpenAuthorizer
	}

	return userInfoProto2Meta(ui)
}
