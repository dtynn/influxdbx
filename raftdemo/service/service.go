package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/dtynn/influxdbx/raftdemo/rpc"
	"github.com/dtynn/influxdbx/raftdemo/rpc/proto"
	"github.com/dtynn/influxdbx/tcp"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	context2 "golang.org/x/net/context"
)

func New(cfg Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

type Service struct {
	cfg Config

	r          *raft.Raft
	raftConfig *raft.Config

	rpcsrv *rpc.ClusterService

	electNotify chan bool
}

func (s *Service) Run(ctx context.Context) error {
	dir, err := filepath.Abs(s.cfg.DataDir)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", s.cfg.Listen)
	if err != nil {
		return err
	}

	mux := tcp.NewMux(ln)

	go func() {
		log.Println("tcp on", ln.Addr().String())
		mux.Serve()
	}()

	s.rpcsrv = rpc.NewCluster(mux.Listen(RPCHeader), s)

	snaps, err := raft.NewFileSnapshotStore(dir, 1, nil)
	if err != nil {
		return err
	}

	storePath := filepath.Join(dir, "raft.db")
	store, err := raftboltdb.NewBoltStore(storePath)
	if err != nil {
		return err
	}

	cacheStore, err := raft.NewLogCache(s.cfg.RaftLogCacheSize, store)
	if err != nil {
		return err
	}

	trans := raft.NewNetworkTransport(&demoStreamLayer{
		Listener: mux.Listen(RAFTHeader),
	}, 5, 5*time.Second, nil)

	s.electNotify = make(chan bool, 1)
	go s.notifyElection()

	s.raftConfig = raft.DefaultConfig()
	s.raftConfig.NotifyCh = s.electNotify
	s.raftConfig.LocalID = raft.ServerID(trans.LocalAddr())

	if len(s.cfg.Join) == 0 {
		hasState, err := raft.HasExistingState(cacheStore, store, snaps)
		if err != nil {
			return err
		}

		if !hasState {
			configuration := raft.Configuration{
				Servers: []raft.Server{
					raft.Server{
						ID:      raft.ServerID(trans.LocalAddr()),
						Address: trans.LocalAddr(),
					},
				},
			}
			if err := raft.BootstrapCluster(s.raftConfig,
				cacheStore, store, snaps, trans, configuration); err != nil {
				return err
			}
		}
	}

	s.r, err = raft.NewRaft(s.raftConfig, &demoFSM{}, cacheStore, store, snaps, trans)
	if err != nil {
		return err
	}

	for _, one := range s.cfg.Join {
		if err := join(one, string(trans.LocalAddr()), string(trans.LocalAddr())); err != nil {
			return err
		}
	}

	return s.rpcsrv.Run(ctx)
}

func (s *Service) notifyElection() {
	for elected := range s.electNotify {
		log.Println("elected", elected)
	}
}

func (s *Service) Join(ctx context2.Context, req *proto.ClusterJoinReq) (*proto.ClusterJoinResp, error) {
	log.Println("try join from", req.GetNode().GetId(), req.GetNode().GetAddress())

	if s.r == nil || s.r.State() != raft.Leader {
		return nil, raft.ErrNotLeader
	}

	f := s.r.AddVoter(raft.ServerID(req.GetNode().GetId()), raft.ServerAddress(req.GetNode().GetAddress()), 0, 0)
	if err := f.Error(); err != nil {
		return nil, err
	}

	return &proto.ClusterJoinResp{}, nil
}

func (s *Service) Add(ctx context2.Context, req *proto.ClusterAddReq) (*proto.ClusterAddResp, error) {
	if s.r == nil {
		return &proto.ClusterAddResp{
			Code: 1,
		}, nil
	}

	log.Println("leader", s.r.Leader())

	f := s.r.Apply([]byte(fmt.Sprintf("%s:%s", req.GetKey(), req.GetValue())), 5*time.Second)
	if err := f.Error(); err != nil {
		log.Println("apply error", err, err == raft.ErrNotLeader)
	}

	return &proto.ClusterAddResp{
		Code: 0,
	}, nil
}

func join(target, id, address string) error {
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		return tcp.Dial("tcp", address, timeout, 11)
	}

	cli, err := rpc.NewClusterClient(target, dialer)
	if err != nil {
		return err
	}

	resp, err := cli.Join(context2.Background(), &proto.ClusterJoinReq{
		Node: &proto.Node{
			Id:      id,
			Address: address,
		},
	})

	if err != nil {
		log.Println("get error", err)
	}

	log.Println("get resp code", resp.GetCode())

	return nil
}
