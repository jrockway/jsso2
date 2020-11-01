//go:generate protoc -I../../protos --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ../../protos/webauthn.proto
package webauthnpb
