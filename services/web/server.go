package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/bredr/go-grpc-example/services/web/graph/generated"
	"github.com/bredr/go-grpc-example/services/web/graph/resolvers"
)

//go:embed www/build
var embeddedFiles embed.FS

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := resolvers.New(context.Background())
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Printf("Unable to process request: %s", err)
		return errors.New("Error processing request")
	})
	srv.Use(apollotracing.Tracer{})

	http.Handle("/", http.FileServer(getFileSystem()))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for app", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getFileSystem() http.FileSystem {

	// Get the build subdirectory as the
	// root directory so that it can be passed
	// to the http.FileServer
	fsys, err := fs.Sub(embeddedFiles, "www/build")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
