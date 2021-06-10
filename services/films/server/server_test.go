package server_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bredr/go-grpc-example/common/client"
	"github.com/bredr/go-grpc-example/proto/films"
	"github.com/bredr/go-grpc-example/services/films/server"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "3.0", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connection := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connection))
		if err != nil {
			return err
		}

		return client.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err = os.Setenv("MONGO_CONNECTION_STRING", connection); err != nil {
		log.Fatalf("Unable to set env variable: %s", err)
	}

	exitStatus := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(exitStatus)
}

func Test_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)
	filmsConn, err := client.DialContext(ctx, "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	filmsClient := films.NewFilmServiceClient(filmsConn)

	result, err := filmsClient.CreateFilm(ctx, &films.Film{Name: "Jaws", Genre: films.FilmGenre_ACTION})
	expect := assert.New(t)
	expect.NoError(err)
	expect.NotNil(result)
}
