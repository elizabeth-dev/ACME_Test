package ports

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type HealthGrpcServer struct {
	dependencies map[string]func(ctx context.Context) error
}

func NewHealthGrpcServer(dependencies map[string]func(ctx context.Context) error) HealthGrpcServer {
	return HealthGrpcServer{dependencies: dependencies}
}

func (s *HealthGrpcServer) Check(
	ctx context.Context,
	in *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	if in.Service == "" {
		for _, pingFunc := range s.dependencies {
			if err := pingFunc(ctx); err != nil {
				return &grpc_health_v1.HealthCheckResponse{
					Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
				}, nil
			}
		}

		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}

	pingFunc, ok := s.dependencies[in.Service]

	if !ok {
		return nil, status.Error(codes.NotFound, "[Check] unknown service")
	}

	if err := pingFunc(ctx); err != nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthGrpcServer) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "[Watch] not implemented")
}
