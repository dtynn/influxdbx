package raft

import (
	"sync"

	"github.com/influxdata/influxdb/uuid"
)

type localQueryManager struct {
	mu     sync.RWMutex
	querys map[uuid.UUID]struct{}
}

func (l *localQueryManager) Check(id uuid.UUID) bool {
	l.mu.RLock()
	_, ok := l.querys[id]
	l.mu.RUnlock()

	return ok
}

func (l *localQueryManager) Add(id uuid.UUID) {
	l.mu.Lock()
	l.querys[id] = struct{}{}
	l.mu.Unlock()
}

func (l *localQueryManager) Done(id uuid.UUID) {
	l.mu.Lock()
	if _, ok := l.querys[id]; ok {
		delete(l.querys, id)
	}
	l.mu.Unlock()
}
