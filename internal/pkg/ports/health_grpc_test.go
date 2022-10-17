package ports

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHealthGrpc(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"gRPC server": {
			"initialize health gRPC server": testNewHealthGrpcServer,
		},
		"check": {
			"call check with no service":     testCheckNoService,
			"call with no service and error": testCheckNoServiceError,
		},
		"check specific service": {
			"call check with service":     testCheckService,
			"call with service and error": testCheckServiceError,
			"call with unknown service":   testCheckUnknownService,
		},
		"watch": {
			"call watcg": testWatch,
		},
	} {
		testGroup := testGroup
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()

				for name, test := range testGroup {
					test := test
					t.Run(
						name, func(t *testing.T) {
							t.Parallel()

							test(t)
						},
					)
				}
			},
		)
	}

}

func testNewHealthGrpcServer(t *testing.T) {
	dependencies := map[string]func(ctx context.Context) error{"test": func(ctx context.Context) error { return nil }}

	out := NewHealthGrpcServer(dependencies)

	assert.NotNil(t, out)
	assert.Equal(t, HealthGrpcServer{dependencies}, out)
}

type dependency struct {
	mock.Mock
}

func (d *dependency) Ping(ctx context.Context) error {
	args := d.Called(ctx)

	return args.Error(0)
}

func testCheckNoService(t *testing.T) {
	mockDependency1 := new(dependency)
	mockDependency2 := new(dependency)
	dependencies := map[string]func(ctx context.Context) error{
		"test1": mockDependency1.Ping,
		"test2": mockDependency2.Ping,
	}

	ctx := context.Background()

	mockDependency1.On("Ping", ctx).Return(nil)
	mockDependency2.On("Ping", ctx).Return(nil)

	srv := NewHealthGrpcServer(dependencies)

	out, err := srv.Check(ctx, &grpc_health_v1.HealthCheckRequest{})

	mockDependency1.AssertNumberOfCalls(t, "Ping", 1)
	mockDependency2.AssertNumberOfCalls(t, "Ping", 1)
	mockDependency1.AssertExpectations(t)
	mockDependency2.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, out.Status)
}

func testCheckNoServiceError(t *testing.T) {
	mockDependency1 := new(dependency)
	mockDependency2 := new(dependency)
	dependencies := map[string]func(ctx context.Context) error{
		"test1": mockDependency1.Ping,
		"test2": mockDependency2.Ping,
	}

	ctx := context.Background()

	mockDependency1.On("Ping", ctx).Return(errors.New("dependency down"))
	mockDependency2.On("Ping", ctx).Return(nil) // This may not be called

	srv := NewHealthGrpcServer(dependencies)

	out, err := srv.Check(ctx, &grpc_health_v1.HealthCheckRequest{})

	mockDependency1.AssertNumberOfCalls(t, "Ping", 1)
	mockDependency1.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, out.Status)
}

func testCheckService(t *testing.T) {
	mockDependency1 := new(dependency)
	mockDependency2 := new(dependency)
	dependencies := map[string]func(ctx context.Context) error{
		"test1": mockDependency1.Ping,
		"test2": mockDependency2.Ping,
	}

	ctx := context.Background()

	mockDependency2.On("Ping", ctx).Return(nil)

	srv := NewHealthGrpcServer(dependencies)

	out, err := srv.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test2"})

	mockDependency1.AssertNumberOfCalls(t, "Ping", 0)
	mockDependency2.AssertNumberOfCalls(t, "Ping", 1)
	mockDependency2.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, out.Status)
}

func testCheckServiceError(t *testing.T) {
	mockDependency1 := new(dependency)
	mockDependency2 := new(dependency)
	dependencies := map[string]func(ctx context.Context) error{
		"test1": mockDependency1.Ping,
		"test2": mockDependency2.Ping,
	}

	ctx := context.Background()

	mockDependency1.On("Ping", ctx).Return(errors.New("dependency down"))

	srv := NewHealthGrpcServer(dependencies)

	out, err := srv.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test1"})

	mockDependency1.AssertNumberOfCalls(t, "Ping", 1)
	mockDependency2.AssertNumberOfCalls(t, "Ping", 0)
	mockDependency1.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, out.Status)
}

func testCheckUnknownService(t *testing.T) {
	mockDependency1 := new(dependency)
	mockDependency2 := new(dependency)
	dependencies := map[string]func(ctx context.Context) error{
		"test1": mockDependency1.Ping,
		"test2": mockDependency2.Ping,
	}

	ctx := context.Background()

	srv := NewHealthGrpcServer(dependencies)

	out, err := srv.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test3"})

	mockDependency1.AssertNumberOfCalls(t, "Ping", 0)
	mockDependency2.AssertNumberOfCalls(t, "Ping", 0)

	assert.Nil(t, out)
	assert.ErrorIs(t, err, status.Error(codes.NotFound, "[Check] unknown service"))
}

type watchServer struct {
	mock.Mock
	grpc.ServerStream
}

func (w *watchServer) Send(resp *grpc_health_v1.HealthCheckResponse) error {
	args := w.Called(resp)

	return args.Error(0)
}

func testWatch(t *testing.T) {
	watchSrv := new(watchServer)
	dependencies := map[string]func(ctx context.Context) error{}

	srv := NewHealthGrpcServer(dependencies)

	err := srv.Watch(&grpc_health_v1.HealthCheckRequest{}, watchSrv)

	watchSrv.AssertNumberOfCalls(t, "Send", 0)

	assert.ErrorIs(t, err, status.Error(codes.Unimplemented, "[Watch] not implemented"))
}
