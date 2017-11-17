package run

import (
	"os"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/services/meta"
)

type Server struct {
	MetaClient *meta.Client

	tcpAddr           string
	reportingDisabled bool

	cfg *Config
}

func NewServer(cfg *Config) (*Server, error) {
	err := os.Mkdir(cfg.Meta.Dir, 0744)
	if err != nil {
		return nil, err
	}

	_, err = influxdb.LoadNode(cfg.Meta.Dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	bind := cfg.BindAddress
	s := &Server{
		MetaClient: meta.NewClient(cfg.Meta),

		tcpAddr:           bind,
		reportingDisabled: cfg.ReportingDisabled,

		cfg: cfg,
	}

	return s, nil
}
