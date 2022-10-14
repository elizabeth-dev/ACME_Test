package main

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/ports"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/service"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/server"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	app := service.NewApplication(ctx)

	server.RunGRPCServer(
		func(server *grpc.Server) {
			svc := ports.NewGrpcServer(app)
			apiV1.RegisterUserServiceServer(server, &svc)
		},
	)
}
