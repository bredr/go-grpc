package server

import (
	"context"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	viper.SetDefault("DEBUG", false)
	viper.AutomaticEnv()

	log.SetLevel(log.WarnLevel)
	if viper.GetBool("DEBUG") {
		log.SetLevel(log.DebugLevel)
	}

}

func New(authFn *grpc_auth.AuthFunc) *grpc.Server {

	if !(authFn != nil) {
		var fn grpc_auth.AuthFunc = func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		}
		authFn = &fn
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(),
		grpc_validator.StreamServerInterceptor(),
		grpc_logrus.StreamServerInterceptor(log.NewEntry(log.StandardLogger())),
		grpc_auth.StreamServerInterceptor(*authFn),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(),
		grpc_validator.UnaryServerInterceptor(),
		grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
		grpc_auth.UnaryServerInterceptor(*authFn),
	}

	server := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				streamInterceptors...,
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				unaryInterceptors...,
			),
		),
	)
	return server
}
