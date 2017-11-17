package run

import (
	"github.com/dtynn/influxdbx/service/raft"
	"github.com/influxdata/influxdb/services/meta"
)

type Config struct {
	Meta *meta.Config `toml:"meta"`
	Raft *raft.Config `toml:"raft"`

	// Server reporting
	ReportingDisabled bool `toml:"reporting-disabled"`

	// BindAddress is the address that all TCP services use (Raft, Snapshot, Cluster, etc.)
	BindAddress string `toml:"bind-address"`
}
