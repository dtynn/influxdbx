package raft

type Config struct {
	LoggingEnabled   bool     `toml:"logging-enabled"`
	LogCacheCapacity int      `toml:"log-cache-capacity"`
	Bootstrap        bool     `toml:"bootstrap"`
	Join             []string `toml:"join"`
	SnapshotRetain   int      `toml:"snapshot-retain"`
	TransportMaxPool int      `toml:"transport-max-pool"`
	TransportTimeout string   `toml:"transport-timeout"`
}

func NewConfig() *Config {
	return &Config{
		LoggingEnabled:   true,
		LogCacheCapacity: 512,
		Bootstrap:        false,
		Join:             []string{},
		SnapshotRetain:   5,
		TransportMaxPool: 5,
		TransportTimeout: "",
	}
}
