package router

import (
	"context"
	"log"

	"github.com/bredr/go-grpc-example/proto/films"
	db "github.com/bredr/go-grpc-example/services/films/repositories/films"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(ctx context.Context) films.FilmServiceServer {
	collection, err := db.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &router{db: collection}
}

type router struct {
	films.UnimplementedFilmServiceServer
	db db.Collection
}

func (r *router) CreateFilm(ctx context.Context, f *films.Film) (*films.Film, error) {
	return r.db.Upsert(ctx, f)
}
func (r *router) UpdateFilm(ctx context.Context, f *films.Film) (*films.Film, error) {
	return r.db.Upsert(ctx, f)
}

func (r *router) FindFilms(ctx context.Context, f *films.FilmSearchRequest) (*films.Films, error) {
	return r.db.Query(ctx, f)
}

func (r *router) DeleteFilm(ctx context.Context, id *films.ID) (*films.Empty, error) {
	if !(id != nil) {
		return nil, status.Error(codes.InvalidArgument, "no id provided")
	}
	err := r.db.Delete(ctx, id.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &films.Empty{}, nil
}
