package main

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/ports"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/service"
	commonPorts "github.com/elizabeth-dev/FACEIT_Test/internal/pkg/ports"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/server"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx := context.Background()

	app, dependencies := service.NewApplication(ctx)

	server.RunGRPCServer(
		func(server *grpc.Server) {
			srv := ports.NewGrpcServer(app)
			healthSrv := commonPorts.NewHealthGrpcServer(dependencies)

			apiV1.RegisterUserServiceServer(server, &srv)
			grpc_health_v1.RegisterHealthServer(server, &healthSrv)
		},
	)
}
