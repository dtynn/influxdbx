package influxdbx

//go:generate protoc -I ./internal/pb --go_out=plugins=grpc:./internal/pb ./internal/pb/common.proto ./internal/pb/meta_data.proto
