module github.com/bredr/go-grpc-example/services/web

go 1.15

replace github.com/bredr/go-grpc-example/proto => ../../proto

replace github.com/bredr/go-grpc-example/common => ../../common

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/bredr/go-grpc-example/common v0.0.0-00010101000000-000000000000
	github.com/bredr/go-grpc-example/proto v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.7.1
	github.com/vektah/gqlparser/v2 v2.1.0
	google.golang.org/protobuf v1.26.0
)
