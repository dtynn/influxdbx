package meta

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/dtynn/influxdbx/internal/pb"
	"github.com/dtynn/influxdbx/service/rpc"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const (
	// MuxHeader tcp mux header
	MuxHeader byte = 6

	raftFile = "raft.db"
)

// NewMeta return meta servie
func NewMeta(dir string, clusterID uint64, cfg *Config) *Meta {
	return &Meta{
		dir:       dir,
		cfg:       cfg,
		clusterID: clusterID,
		logger:    zap.NewNop(),
	}
}

// Meta meta manager service
type Meta struct {
	dir string
	cfg *Config

	clusterID uint64

	mu sync.RWMutex

	logger *zap.Logger

	net.Listener

	r         *raft.Raft
	rconf     *raft.Config
	raftstore *raftboltdb.BoltStore

	electNotify chan bool
}

// Open open a meta service
func (m *Meta) Open() error {
	path := filepath.Join(m.dir, raftFile)

	store, err := raftboltdb.NewBoltStore(path)
	if err != nil {
		return err
	}

	logs, err := raft.NewLogCache(m.cfg.LogCacheCapacity, store)
	if err != nil {
		return err
	}

	snaps, err := raft.NewFileSnapshotStore(m.dir, m.cfg.SnapshotRetain, nil)
	if err != nil {
		return err
	}

	fsm := newFSM()

	trans := raft.NewNetworkTransport(newStream(m.Listener), 5, 10*time.Second, nil)

	m.electNotify = make(chan bool, 1)

	localID := raft.ServerID(strconv.FormatUint(m.clusterID, 10))

	m.rconf = raft.DefaultConfig()
	m.rconf.LocalID = localID
	m.rconf.NotifyCh = m.electNotify

	local := raft.Server{
		ID:      localID,
		Address: trans.LocalAddr(),
	}

	if m.cfg.Bootstrap {
		hasState, err := raft.HasExistingState(logs, store, snaps)
		if err != nil {
			return err
		}

		if !hasState {
			configuration := raft.Configuration{}
			configuration.Servers = append(configuration.Servers, local)

			if err := raft.BootstrapCluster(m.rconf, logs, store, snaps, trans, configuration); err != nil {
				return err
			}
		}
	}

	m.r, err = raft.NewRaft(m.rconf, fsm, logs, store, snaps, trans)
	if err != nil {
		return err
	}

	if m.cfg.Join != "" {
		if err := join(m.cfg.Join, m.clusterID, local.Address); err != nil {
			return err
		}
	}

	go m.onElect()

	return nil
}

func (m *Meta) onElect() {
	for elected := range m.electNotify {
		m.logger.Info(fmt.Sprintf("elected as leader %t", elected))
	}
}

// Close close the meta service
func (m *Meta) Close() error {
	if m.r != nil {
		f := m.r.Shutdown()
		f.Error()
	}

	if m.raftstore != nil {
		m.raftstore.Close()
	}

	if m.electNotify != nil {
		if m.rconf != nil {
			m.rconf.NotifyCh = nil
		}

		close(m.electNotify)
	}

	return nil
}

// WithLogger setup new logger
func (m *Meta) WithLogger(l *zap.Logger) {
	m.logger = l.With(zap.String("service", "meta"))
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

func join(target string, localID uint64, localAddress raft.ServerAddress) error {
	cli := rpc.NewMetaClient(target, true)
	resp, err := cli.Join(context.Background(), &pb.MetaJoinReq{
		Id:      localID,
		Address: string(localAddress),
	})

	if err != nil {
		return err
	}

	if resp.GetResult().GetCode() != pb.ResultCode_ResultCodeOK {
		return fmt.Errorf("join result: %#v", resp.GetResult())
	}

	return nil
}
