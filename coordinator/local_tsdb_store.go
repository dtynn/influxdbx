package coordinator

import "github.com/influxdata/influxdb/tsdb"

var _ TSDBStore = (*LocalTSDBStore)(nil)

// LocalTSDBStore embeds a tsdb.Store and implements IteratorCreator
// to satisfy the TSDBStore interface.
type LocalTSDBStore struct {
	*tsdb.Store
}

func (l *LocalTSDBStore) ShardGroup(ids []uint64) (tsdb.ShardGroup, error) {
	return l.Store.ShardGroup(ids), nil
}

func (l *LocalTSDBStore) ShardIDs() ([]uint64, error) {
	return l.Store.ShardIDs(), nil
}
