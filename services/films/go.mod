module github.com/bredr/go-grpc-example/services/films

go 1.16

replace github.com/bredr/go-grpc-example/proto => ../../proto

replace github.com/bredr/go-grpc-example/common => ../../common

require (
	github.com/bredr/go-grpc-example/common v0.0.0-00010101000000-000000000000
	github.com/bredr/go-grpc-example/proto v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/ory/dockertest/v3 v3.6.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1 // indirect
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb // indirect
	golang.org/x/oauth2 v0.0.0-20201203001011-0b49973bad19 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201203001206-6486ece9c497 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/grpc/examples v0.0.0-20210511173931-5f95ad62331a // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
