package main

//go:generate protoc -I ./service/raft/internal --go_out=plugins=grpc:./service/raft/internal ./service/raft/internal/cluster.proto
