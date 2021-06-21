package server_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bredr/go-grpc-example/common/client"
	"github.com/bredr/go-grpc-example/proto/films"
	"github.com/bredr/go-grpc-example/services/films/server"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerTestSuite struct {
	suite.Suite
	mgoClient *mongo.Client
	pool      *dockertest.Pool
	resource  *dockertest.Resource
}

func (suite *ServerTestSuite) SetupSuite() {
	var err error
	suite.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	suite.resource, err = suite.pool.Run("mongo", "3.0", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connection := fmt.Sprintf("mongodb://localhost:%s", suite.resource.GetPort("27017/tcp"))
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := suite.pool.Retry(func() error {
		var err error
		suite.mgoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connection))
		if err != nil {
			return err
		}

		return suite.mgoClient.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err = os.Setenv("MONGO_CONNECTION_STRING", connection); err != nil {
		log.Fatalf("Unable to set env variable: %s", err)
	}
}

func (suite *ServerTestSuite) TearDownSuite() {
	if err := suite.pool.Purge(suite.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (suite *ServerTestSuite) SetupTest() {
	if err := suite.mgoClient.Database("example").Drop(context.Background()); err != nil {
		log.Fatalf("Could not drop database: %s", err)
	}
}

func (suite *ServerTestSuite) TestCreate() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	result, err := filmsClient.CreateFilm(ctx, &films.Film{Name: "Jaws", Genre: films.FilmGenre_ACTION})
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *ServerTestSuite) serverClient() (films.FilmServiceClient, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx)
	filmsConn, err := client.DialContext(ctx, "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	filmsClient := films.NewFilmServiceClient(filmsConn)
	return filmsClient, ctx, cancel
}

func (suite *ServerTestSuite) TestFindFilms() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	_, err := filmsClient.CreateFilm(ctx, &films.Film{Name: "Jaws", Genre: films.FilmGenre_ACTION})
	suite.NoError(err)
	result, err := filmsClient.FindFilms(ctx, &films.FilmSearchRequest{})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 1)
}

func (suite *ServerTestSuite) TestFindFilmsByGenre() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	_, err := filmsClient.CreateFilm(ctx, &films.Film{Name: "Jaws", Genre: films.FilmGenre_ACTION})
	suite.NoError(err)
	result, err := filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		AllowedGenres: []films.FilmGenre{films.FilmGenre_ACTION},
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 1)
	result, err = filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		AllowedGenres: []films.FilmGenre{films.FilmGenre_ROMANCE},
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 0)
}

func (suite *ServerTestSuite) TestFindFilmsBySearchTerm() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	_, err := filmsClient.CreateFilm(ctx, &films.Film{Name: "Jaws", Genre: films.FilmGenre_ACTION})
	suite.NoError(err)
	result, err := filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		NameSearch: "Jaw",
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 1)
	result, err = filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		NameSearch: "Cat",
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 0)
}

func (suite *ServerTestSuite) TestFindFilmsByReleaseDate() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	_, err := filmsClient.CreateFilm(ctx, &films.Film{
		Name: "Jaws", Genre: films.FilmGenre_ACTION,
		ReleaseDate: timestamppb.New(time.Date(2000, time.January, 2, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	result, err := filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		ReleasedAfter: timestamppb.New(time.Date(2000, time.January, 1, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 1)
	result, err = filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		ReleasedAfter: timestamppb.New(time.Date(2000, time.January, 3, 11, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 0)
}

func (suite *ServerTestSuite) TestFindFilmsByAllFilters() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	_, err := filmsClient.CreateFilm(ctx, &films.Film{
		Name: "Jaws", Genre: films.FilmGenre_ACTION,
		ReleaseDate: timestamppb.New(time.Date(2000, time.January, 2, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	result, err := filmsClient.FindFilms(ctx, &films.FilmSearchRequest{
		NameSearch:    "Jaw",
		AllowedGenres: []films.FilmGenre{films.FilmGenre_ACTION},
		ReleasedAfter: timestamppb.New(time.Date(2000, time.January, 1, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	suite.Len(result.GetFilms(), 1)
}

func (suite *ServerTestSuite) TestDeleteFilm() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	film, err := filmsClient.CreateFilm(ctx, &films.Film{
		Name: "Jaws", Genre: films.FilmGenre_ACTION,
		ReleaseDate: timestamppb.New(time.Date(2000, time.January, 2, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	_, err = filmsClient.DeleteFilm(ctx, &films.ID{ID: film.ID})
	suite.NoError(err)
	count, err := suite.mgoClient.Database("example").Collection("films").CountDocuments(ctx, bson.D{})
	suite.NoError(err)
	suite.EqualValues(0, count)
}

func (suite *ServerTestSuite) TestUpdateFilm() {
	filmsClient, ctx, cancel := suite.serverClient()
	defer cancel()
	film, err := filmsClient.CreateFilm(ctx, &films.Film{
		Name: "Jaws", Genre: films.FilmGenre_ACTION,
		ReleaseDate: timestamppb.New(time.Date(2000, time.January, 2, 12, 03, 02, 0, time.UTC)),
	})
	suite.NoError(err)
	film.Genre = films.FilmGenre_ROMANCE
	updatedFilm, err := filmsClient.UpdateFilm(ctx, film)
	suite.NoError(err)
	suite.Equal(films.FilmGenre_ROMANCE, updatedFilm.Genre)
	count, err := suite.mgoClient.Database("example").Collection("films").CountDocuments(ctx, bson.D{})
	suite.NoError(err)
	suite.EqualValues(1, count)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
