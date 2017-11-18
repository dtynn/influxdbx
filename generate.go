package influxdbx

//go:generate protoc -I ./cluster/proto --go_out=plugins=grpc:./cluster/proto ./cluster/proto/cluster.proto
//go:generate protoc -I ./service/raft/internal --go_out=plugins=grpc:./service/raft/internal ./service/raft/internal/cmd.proto
