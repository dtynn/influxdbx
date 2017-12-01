package rpc

// Config rpc service config
type Config struct {
	EnableLog bool `toml:"enable-log"`
}

// NewConfig new config
func NewConfig() Config {
	return Config{
		EnableLog: true,
	}
}
