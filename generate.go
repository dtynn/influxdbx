package influxdbx

//go:generate protoc -I ./cluster/proto --go_out=plugins=grpc:./cluster/proto ./cluster/proto/common.proto ./cluster/proto/cluster.proto ./cluster/proto/meta.proto ./cluster/proto/store.proto
//go:generate protoc -I ./raft/internal --go_out=plugins=grpc:./raft/internal ./raft/internal/cmd.proto
