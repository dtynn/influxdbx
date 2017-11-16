package service

type Config struct {
	Listen           string
	RPCAddr          string
	RaftAddr         string
	DataDir          string
	Join             []string
	RaftLogCacheSize int
}
