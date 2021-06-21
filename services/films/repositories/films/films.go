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
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(ctx context.Context, h *health.Server) (Collection, error) {
	viper.SetDefault("MONGO_CONNECTION_STRING", "mongodb://localhost:27017")
	viper.AutomaticEnv()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGO_CONNECTION_STRING")))
	if err != nil {
		return nil, err
	}

	c := client.Database("example").Collection("films")

	_, err = c.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "name", Value: "text"}}})
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Ping(ctx, nil); err != nil {
					h.SetServingStatus("mongo", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
				} else {
					h.SetServingStatus("mongo", grpc_health_v1.HealthCheckResponse_SERVING)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

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
	filters := []bson.M{}
	if query != nil {
		if len(query.AllowedGenres) > 0 {
			genres := make([]bson.M, len(query.AllowedGenres))
			for i, x := range query.AllowedGenres {
				genres[i] = bson.M{"genre": x.String()}
			}
			filters = append(filters, bson.M{"$or": genres})
		}
		if query.ReleasedAfter != nil {
			filters = append(filters, bson.M{"release_date": bson.M{"$gte": query.ReleasedAfter.AsTime()}})
		}
		if query.NameSearch != "" {
			filters = append(filters, bson.M{"$text": bson.M{"$search": query.NameSearch}})
		}
	}
	var filter bson.M
	if len(filters) == 1 {
		filter = filters[0]
	}
	if len(filters) > 1 {
		filter = bson.M{"$and": filters}
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
