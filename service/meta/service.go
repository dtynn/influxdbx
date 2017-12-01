package meta

import (
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

// Meta meta manager
type Meta struct {
	clusterID uint64

	mu sync.RWMutex

	r      *raft.Raft
	logger *zap.Logger
}

// Open open a meta service
func (m *Meta) Open() error {
	return nil
}

// WithLogger setup new logger
func (m *Meta) WithLogger(l *zap.Logger) {
	m.logger = l.With(zap.String("service", "meta"))
}

// Close close meta service
func (m *Meta) Close() error {
	return nil
}

// AddMetaNode add meta node to raft
func (m *Meta) AddMetaNode(id uint64, address string) error {

	return m.r.AddVoter(
		raft.ServerID(strconv.FormatUint(id, 10)),
		raft.ServerAddress(address),
		0,
		5*time.Second,
	).Error()
}

// RemoveMetaNode remove meta node from raft
func (m *Meta) RemoveMetaNode(id uint64) error {

	return m.r.RemoveServer(raft.ServerID(strconv.FormatUint(id, 10)), 0, 5*time.Second).Error()
}

// IsLeader if current node is leader
func (m *Meta) IsLeader() bool {
	return m.r.State() == raft.Leader
}

// Leader return address of leader node
func (m *Meta) Leader() string {
	return string(m.r.Leader())
}
