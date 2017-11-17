package main

//go:generate protoc -I ./service/cluster/proto --go_out=plugins=grpc:./service/cluster/proto ./service/cluster/proto/cluster.proto
//go:generate protoc -I ./service/raft/internal --go_out=plugins=grpc:./service/raft/internal ./service/raft/internal/cmd.proto
