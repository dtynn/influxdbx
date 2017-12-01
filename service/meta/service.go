package meta

import (
	"sync"

	"github.com/dtynn/influxdbx/internal/pb"
	"github.com/hashicorp/raft"
	"golang.org/x/net/context"
)

var (
	_ pb.MetaServer = (*Meta)(nil)
)

// Meta meta manager
type Meta struct {
	clusterID uint64

	mu sync.RWMutex

	r *raft.Raft
}

// Open open a meta service
func (m *Meta) Open() error {
	return nil
}

// Join api for meta service
func (m *Meta) Join(ctx context.Context, req *pb.MetaJoinReq) (*pb.MetaJoinResp, error) {
	return nil, nil
}

// Remove api for meta service
func (m *Meta) Remove(ctx context.Context, req *pb.MetaRemoveReq) (*pb.MetaRemoveResp, error) {
	return nil, nil
}

// Apply api for meta service
func (m *Meta) Apply(ctx context.Context, req *pb.MetaApplyReq) (*pb.MetaApplyResp, error) {
	return nil, nil
}
