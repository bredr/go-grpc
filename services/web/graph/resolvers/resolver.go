package resolvers

import (
	"context"
	"log"

	"github.com/bredr/go-grpc-example/common/client"
	"github.com/spf13/viper"

	"github.com/bredr/go-grpc-example/proto/films"
)

//go:generate rm -rf ../generated
//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
const FILMS_URL = "FILMS_URL"

type Resolver struct {
	films films.FilmServiceClient
}

func New(ctx context.Context) *Resolver {
	viper.SetDefault(FILMS_URL, "films")
	viper.AutomaticEnv()
	filmsConn, err := client.DialContext(ctx, viper.GetString(FILMS_URL))
	if err != nil {
		log.Fatal(err)
	}
	filmsClient := films.NewFilmServiceClient(filmsConn)
	return &Resolver{films: filmsClient}
}
