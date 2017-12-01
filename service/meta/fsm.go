package meta

import (
	"fmt"
	"io"
	"sync"

	"github.com/dtynn/influxdbx/internal/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/raft"
)

var (
	_ raft.FSM         = (*FSM)(nil)
	_ raft.FSMSnapshot = (*FSMSnapshot)(nil)
)

func newFSM() *FSM {
	return &FSM{}
}

// FSM implement raft.FSM
type FSM struct {
	mu   sync.RWMutex
	data *Data
}

// Apply apply raft.Log
func (f *FSM) Apply(log *raft.Log) interface{} {
	return fmt.Errorf("not implemented")
}

// Snapshot return raft.Snapshot
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, fmt.Errorf("not implemented")
}

// Restore restore data from reader
func (f *FSM) Restore(r io.ReadCloser) error {
	return fmt.Errorf("not impl")
}

func (f *FSM) applyCreateMetaNode(buf []byte) error {
	var cmd pb.MetaCmdCreateMetaNode
	if err := proto.Unmarshal(buf, &cmd); err != nil {
		return err
	}

	return nil
}

// FSMSnapshot implement raft.FSMSnapshot
type FSMSnapshot struct {
}

// Persist persist local data
func (f *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return fmt.Errorf("not impl")
}

// Release release snapshot
func (f *FSMSnapshot) Release() {

}
