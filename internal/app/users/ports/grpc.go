package ports

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/app"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/app/command"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/app/query"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/grpc_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g *GrpcServer) CreateUser(ctx context.Context, request *apiV1.CreateUserRequest) (*apiV1.User, error) {
	cmd := command.CreateUser{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Nickname:  request.GetNickname(),
		Password:  request.GetPassword(),
		Email:     request.GetEmail(),
		Country:   request.GetCountry(),
	}

	id, err := g.app.Commands.CreateUser.Handle(ctx, cmd)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newUser, err := g.app.Queries.GetUserById.Handle(ctx, id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiV1.User{
		Id:        newUser.Id,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Nickname:  newUser.Nickname,
		Password:  newUser.Password,
		Email:     newUser.Email,
		Country:   newUser.Country,
		CreatedAt: timestamppb.New(newUser.CreatedAt),
		UpdatedAt: timestamppb.New(newUser.UpdatedAt),
	}, nil
}

func (g *GrpcServer) GetUsers(request *apiV1.GetUsersRequest, srv apiV1.UserService_GetUsersServer) error {
	var filters []query_utils.Filter
	for _, filter := range request.GetFilters() {
		filters = append(filters, grpc_utils.MapGrpcFilterToFilter(filter))
	}

	var sorts []query_utils.Sort
	for _, sort := range request.GetSort() {
		sorts = append(sorts, grpc_utils.MapGrpcSortToSort(sort))
	}

	getUsersQuery := query.GetUsers{
		Filters: filters,
		Sort:    sorts,
		Pagination: query_utils.Pagination{
			Limit:  request.GetPagination().GetLimit(),
			Offset: request.GetPagination().GetOffset(),
		},
	}

	users, err := g.app.Queries.GetUsers.Handle(srv.Context(), getUsersQuery)

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, user := range users {
		if err := srv.Send(
			&apiV1.User{
				Id:        user.Id,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Nickname:  user.Nickname,
				Password:  user.Password,
				Email:     user.Email,
				Country:   user.Country,
				CreatedAt: timestamppb.New(user.CreatedAt),
				UpdatedAt: timestamppb.New(user.UpdatedAt),
			},
		); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (g *GrpcServer) UpdateUser(ctx context.Context, request *apiV1.UpdateUserRequest) (*apiV1.User, error) {

	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "[UpdateUser] id is required")
	}

	cmd := command.UpdateUser{
		Id:        request.GetId(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Nickname:  request.Nickname,
		Password:  request.Password,
		Email:     request.Email,
		Country:   request.Country,
	}

	err := g.app.Commands.UpdateUser.Handle(ctx, cmd)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	updatedUser, err := g.app.Queries.GetUserById.Handle(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiV1.User{
		Id:        updatedUser.Id,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Nickname:  updatedUser.Nickname,
		Password:  updatedUser.Password,
		Email:     updatedUser.Email,
		Country:   updatedUser.Country,
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
		UpdatedAt: timestamppb.New(updatedUser.UpdatedAt),
	}, nil
}

func (g *GrpcServer) RemoveUser(ctx context.Context, request *apiV1.RemoveUserRequest) (*emptypb.Empty, error) {
	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "[Removeuser] id is required")
	}

	err := g.app.Commands.RemoveUser.Handle(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}