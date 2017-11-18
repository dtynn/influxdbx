package run

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/dtynn/influxdbx/cluster"
	"github.com/dtynn/influxdbx/service/raft"
	"github.com/influxdata/influxdb/services/meta"
)

const (
	DefaultBindAddress = "0.0.0.0:8117"
)

func NewConfig() *Config {
	return &Config{
		Meta:    meta.NewConfig(),
		Raft:    raft.NewConfig(),
		Cluster: cluster.NewConfig(),

		ReportingDisabled: false,
		BindAddress:       DefaultBindAddress,
	}
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()

	var homeDir string
	// By default, store meta and data files in current users home directory
	u, err := user.Current()
	if err == nil {
		homeDir = u.HomeDir
	} else if os.Getenv("HOME") != "" {
		homeDir = os.Getenv("HOME")
	} else {
		return nil, fmt.Errorf("failed to determine current user for storage")
	}

	c.Meta.Dir = filepath.Join(homeDir, ".influxdbx/meta")
	c.Raft.Bootstrap = true

	return c, nil
}

type Config struct {
	Meta    *meta.Config    `toml:"meta"`
	Raft    *raft.Config    `toml:"raft"`
	Cluster *cluster.Config `toml:"cluster"`

	// Server reporting
	ReportingDisabled bool `toml:"reporting-disabled"`

	// BindAddress is the address that all TCP services use (Raft, Snapshot, Cluster, etc.)
	BindAddress string `toml:"bind-address"`
}

func (c *Config) Load(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}
