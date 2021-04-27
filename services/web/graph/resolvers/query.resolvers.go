package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/bredr/go-grpc-example/services/web/graph/generated"
	"github.com/bredr/go-grpc-example/services/web/graph/model"
	"github.com/bredr/go-grpc-example/services/web/services/films"
)

func (r *filmResolver) Actors(ctx context.Context, obj *model.Film) ([]model.Actor, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Films(ctx context.Context, search *model.FilmSearch) ([]model.Film, error) {
	result, err := r.films.FindFilms(ctx, films.GenerateSearchRequest(search))
	if err != nil {
		return nil, err
	}
	return films.MapFilms(result.GetFilms()), nil
}

// Film returns generated.FilmResolver implementation.
func (r *Resolver) Film() generated.FilmResolver { return &filmResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type filmResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
