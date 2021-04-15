//go:generate protoc --proto_path=.. --go_out=:.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=.. web/web.proto
package web
