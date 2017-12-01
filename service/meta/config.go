package meta

const (
	// DefaultLogCacheCapacity default capacity for raft.LogCache
	DefaultLogCacheCapacity = 128

	// DefaultSnapshotRetain default snapshot retain for raft.FileSnapshotStore
	DefaultSnapshotRetain = 1
)

// Config meta service config
type Config struct {
	LogCacheCapacity int    `toml:"log-cache-capacity"`
	SnapshotRetain   int    `toml:"snapshot-retain"`
	Bootstrap        bool   `toml:"bootstrap"`
	Join             string `toml:"join"`
}

// NewConfig return new config
func NewConfig() Config {
	return Config{
		LogCacheCapacity: DefaultLogCacheCapacity,
		SnapshotRetain:   DefaultSnapshotRetain,
	}
}
