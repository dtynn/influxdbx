package run

import (
	"github.com/influxdata/influxdb/services/meta"
)

type Config struct {
	Meta *meta.Config `toml:"meta"`

	// Server reporting
	ReportingDisabled bool `toml:"reporting-disabled"`

	// BindAddress is the address that all TCP services use (Raft, Snapshot, Cluster, etc.)
	BindAddress string `toml:"bind-address"`
}
