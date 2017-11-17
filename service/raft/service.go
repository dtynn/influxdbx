package raft

import (
	"context"
	"net"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	MuxHeader byte = 11

	raftFile = "raft.db"

	defaultTransportTimeout = 10 * time.Second
)

type Raft struct {
	dir          string
	cfg          *Config
	transTimeout time.Duration

	r     *raft.Raft
	rconf *raft.Config

	Listner net.Listener

	store *raftboltdb.BoltStore

	electNotify chan bool
}

func (r *Raft) Open(ctx context.Context) error {
	defer func() {
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
	}()

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

	trans := raft.NewNetworkTransport(NewStream(r.Listner), r.cfg.TransportMaxPool, r.transTimeout, nil)
	if err != nil {
		return err
	}

	r.electNotify = make(chan bool, 1)
	go r.onElect()

	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(trans.LocalAddr())

	if r.cfg.Bootstrap {
		hasState, err := raft.HasExistingState(logs, store, snaps)
		if err != nil {
			return err
		}

		if !hasState {
			configuration := raft.Configuration{}
			configuration.Servers = append(configuration.Servers, raft.Server{
				ID:      raft.ServerID(trans.LocalAddr()),
				Address: trans.LocalAddr(),
			})

			if err := raft.BootstrapCluster(conf, logs, store, snaps, trans, configuration); err != nil {
				return err
			}
		}
	}

	r.r, err = raft.NewRaft(conf, fsm, logs, store, snaps, trans)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (r *Raft) onElect() {
	for _ = range r.electNotify {

	}
}
