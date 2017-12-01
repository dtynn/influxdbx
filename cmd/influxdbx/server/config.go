package server

import (
	"github.com/BurntSushi/toml"
	"github.com/dtynn/influxdbx/service/meta"
	"github.com/dtynn/influxdbx/service/rpc"
)

const (
	// DefaultBindAddress default bind address for server
	DefaultBindAddress = "0.0.0.0:8117"
)

// NewConfig new default config
func NewConfig() *Config {
	return &Config{
		Meta: meta.NewConfig(),
		RPC:  rpc.NewConfig(),

		ReportingDisabled: false,
		BindAddress:       DefaultBindAddress,
	}
}

// Config server config
type Config struct {
	Meta *meta.Config `toml:"meta"`
	RPC  *rpc.Config  `toml:"rpc"`

	Dir string `toml:"dir"`

	// Server reporting
	ReportingDisabled bool `toml:"reporting-disabled"`

	// BindAddress is the address that all TCP services use (Raft, Snapshot, Cluster, etc.)
	BindAddress string `toml:"bind-address"`
}

// Load load from file
func (c *Config) Load(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}
