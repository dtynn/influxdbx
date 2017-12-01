package rpc

// Config rpc service config
type Config struct {
	Verbose bool `toml:"verbose"`
}

// NewConfig new config
func NewConfig() *Config {
	return &Config{
		Verbose: true,
	}
}
