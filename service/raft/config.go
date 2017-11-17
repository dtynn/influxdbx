package raft

type Config struct {
	LogCacheCapacity int      `toml:"log-cache-capacity"`
	Bootstrap        bool     `toml:"bootstrap"`
	Join             []string `toml:"join"`
	SnapshotRetain   int      `toml:"snapshot-retain"`
	TransportMaxPool int      `toml:"transport-max-pool"`
	TransportTimeout string   `toml:"transport-timeout"`
}
