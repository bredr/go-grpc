package films

import (
	"context"
	"time"

	pb "github.com/bredr/go-grpc-example/proto/films"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(ctx context.Context) (Collection, error) {
	viper.SetDefault("MONGO_CONNECTION_STRING", "mongodb://localhost:27017")
	viper.AutomaticEnv()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGO_CONNECTION_STRING")))
	if err != nil {
		return nil, err
	}

	c := client.Database("example").Collection("films")

	go func() {
		<-ctx.Done()
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return &collection{c}, nil
}

type Film struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	ReleaseDate time.Time          `bson:"release_date,omitempty"`
	Genre       string             `bson:"genre,omitempty"`
	ActorIDs    []string           `bson:"actors,omitempty"`
}

type Collection interface {
	Upsert(ctx context.Context, film *pb.Film) (*pb.Film, error)
	Delete(ctx context.Context, id string) error
	Query(ctx context.Context, query *pb.FilmSearchRequest) (*pb.Films, error)
}

type collection struct {
	c *mongo.Collection
}

func (c *collection) Upsert(ctx context.Context, f *pb.Film) (*pb.Film, error) {
	if !(f != nil) {
		return nil, nil
	}
	if f.ID == "" {
		f.ID = primitive.NewObjectID().Hex()
	}
	id, err := primitive.ObjectIDFromHex(f.ID)
	if err != nil {
		return nil, err
	}
	film := Film{
		ID:          id,
		Name:        f.Name,
		Genre:       f.GetGenre().String(),
		ReleaseDate: f.ReleaseDate.AsTime(),
		ActorIDs:    []string{},
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": film.ID}
	if _, err := c.c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: film}}, opts); err != nil {
		return nil, err
	}
	return f, nil
}

func (c *collection) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.c.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (c *collection) Query(ctx context.Context, query *pb.FilmSearchRequest) (*pb.Films, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: -1}})
	filter := bson.D{}
	if query != nil {
		if len(query.AllowedGenres) > 0 {
			genres := make([]bson.E, len(query.AllowedGenres))
			for i, x := range query.AllowedGenres {
				genres[i] = bson.E{Key: "genre", Value: x}
			}
			filter = append(filter, bson.E{Key: "$or", Value: genres})
		}
		if query.ReleasedAfter != nil {
			filter = append(filter, bson.E{Key: "release_date", Value: bson.E{Key: "$ge", Value: query.ReleasedAfter.AsTime()}})
		}
		if query.NameSearch != "" {
			filter = append(filter, bson.E{Key: "name", Value: bson.E{Key: "$text", Value: bson.E{Key: "$search", Value: query.NameSearch}}})
		}
	}
	cursor, err := c.c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []Film
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	parsedResults := make([]*pb.Film, len(results))
	for i, x := range results {
		genre := pb.FilmGenre(pb.FilmGenre_value[x.Genre])
		parsedResults[i] = &pb.Film{
			ID:          x.ID.Hex(),
			Name:        x.Name,
			Genre:       genre,
			ReleaseDate: timestamppb.New(x.ReleaseDate),
		}
	}
	return &pb.Films{Films: parsedResults}, nil
}
