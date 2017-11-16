package service

import (
	"io"
	"log"

	"github.com/hashicorp/raft"
)

var (
	_ raft.FSM         = &demoFSM{}
	_ raft.FSMSnapshot = &demoFSMSnapshot{}
)

type demoFSM struct {
}

func (d *demoFSM) Apply(l *raft.Log) interface{} {
	log.Println("FSM.Apply", string(l.Data))
	return string(l.Data)
}

func (d *demoFSM) Snapshot() (raft.FSMSnapshot, error) {
	log.Println("FSM.Snapshot")
	return &demoFSMSnapshot{}, nil
}

func (d *demoFSM) Restore(rc io.ReadCloser) error {
	defer rc.Close()

	log.Println("FSM.Restore")

	return nil
}

type demoFSMSnapshot struct {
}

func (d *demoFSMSnapshot) Persist(sink raft.SnapshotSink) error {
	log.Println("FSMSnapshot.Persist")
	return nil
}

func (d *demoFSMSnapshot) Release() {
	log.Println("FSMSnapshot.Release")
}
