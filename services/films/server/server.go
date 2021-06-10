package server

import (
	"context"
	"log"
	"net"

	"github.com/bredr/go-grpc-example/common/server"
	"github.com/bredr/go-grpc-example/proto/films"
	"github.com/bredr/go-grpc-example/services/films/router"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) {

	viper.SetDefault("PORT", "3000")
	viper.AutomaticEnv()
	srv := server.New(nil)

	films.RegisterFilmServiceServer(srv, router.New(ctx))

	lis, err := net.Listen("tcp", "0.0.0.0:"+viper.GetString("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	srv.Stop()
}
