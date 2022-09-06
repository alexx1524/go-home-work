//nolint:lll
//go:generate sh -c "protoc --proto_path=. *.proto --go_out ../internal/server/grpc/eventpb --go-grpc_out ../internal/server/grpc/eventpb"
package api
