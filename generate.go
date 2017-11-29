package influxdbx

//go:generate protoc -I ./internal/proto --go_out=plugins=grpc:./internal/proto ./internal/proto/common.proto ./internal/proto/meta.proto
