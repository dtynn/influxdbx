package main

//go:generate protoc -I ./rpc/proto --go_out=plugins=grpc:./rpc/proto ./rpc/proto/cluster.proto
