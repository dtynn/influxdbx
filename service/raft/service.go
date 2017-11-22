package raft

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/dtynn/influxdbx/cluster"
	"github.com/dtynn/influxdbx/cluster/proto"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/influxdata/influxdb/services/meta"
	"go.uber.org/zap"
)

const (
	MuxHeader byte = 12

	raftFile = "raft.db"

	defaultTransportTimeout = 10 * time.Second
)

func NewRaft(dir string, cfg *Config) *Raft {
	transTimeout, _ := time.ParseDuration(cfg.TransportTimeout)
	if transTimeout <= 0 {
		transTimeout = defaultTransportTimeout
	}

	return &Raft{
		dir:          dir,
		cfg:          cfg,
		transTimeout: transTimeout,
		logger:       zap.NewNop(),
	}
}

type Raft struct {
	dir          string
	cfg          *Config
	transTimeout time.Duration

	logger *zap.Logger

	net.Listener

	r     *raft.Raft
	rconf *raft.Config

	store *raftboltdb.BoltStore

	electNotify chan bool

	MetaClient *meta.Client
}

func (r *Raft) Open() error {
	path := filepath.Join(r.dir, raftFile)

	store, err := raftboltdb.NewBoltStore(path)
	if err != nil {
		return err
	}

	logs, err := raft.NewLogCache(r.cfg.LogCacheCapacity, store)
	if err != nil {
		return err
	}

	snaps, err := raft.NewFileSnapshotStore(r.dir, r.cfg.SnapshotRetain, nil)
	if err != nil {
		return err
	}

	fsm := &FSM{}

	trans := raft.NewNetworkTransport(NewStream(r.Listener), r.cfg.TransportMaxPool, r.transTimeout, nil)
	if err != nil {
		return err
	}

	r.electNotify = make(chan bool, 1)

	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(trans.LocalAddr())
	conf.NotifyCh = r.electNotify

	local := raft.Server{
		ID:      raft.ServerID(trans.LocalAddr()),
		Address: trans.LocalAddr(),
	}

	if r.cfg.Bootstrap {
		hasState, err := raft.HasExistingState(logs, store, snaps)
		if err != nil {
			return err
		}

		if !hasState {
			configuration := raft.Configuration{}
			configuration.Servers = append(configuration.Servers, local)

			if err := raft.BootstrapCluster(conf, logs, store, snaps, trans, configuration); err != nil {
				return err
			}
		}
	}

	r.r, err = raft.NewRaft(conf, fsm, logs, store, snaps, trans)
	if err != nil {
		return err
	}

	for _, target := range r.cfg.Join {
		if err := join(target, local.ID, local.Address); err != nil {
			return err
		}
	}

	go r.onElect()

	return nil
}

func (r *Raft) Close() error {
	if r.r != nil {
		f := r.r.Shutdown()
		f.Error()
	}

	if r.store != nil {
		r.store.Close()
	}

	if r.electNotify != nil {
		if r.rconf != nil {
			r.rconf.NotifyCh = nil
		}

		close(r.electNotify)
	}

	return nil
}

func (r *Raft) WithLogger(l *zap.Logger) {
	r.logger = l.With(zap.String("service", "raft"))
}

func (r *Raft) onElect() {
	for elected := range r.electNotify {
		r.logger.Info(fmt.Sprintf("raft election won %v", elected))
	}
}

func join(target string, id raft.ServerID, address raft.ServerAddress) error {
	_, err := cluster.NewClusterClient(target, true).Join(context.Background(), &proto.ClusterJoinReq{
		Node: []*proto.Node{
			{
				Id:      string(id),
				Address: string(address),
			},
		},
	})

	return err
}
