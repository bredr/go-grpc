package main

import (
	"context"

	"github.com/bredr/go-grpc-example/services/films/server"
)

func main() {
	server.Run(context.Background())
}
