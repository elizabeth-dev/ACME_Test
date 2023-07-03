package ports

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/command"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/query"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/errors"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/grpc_utils"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils"
	apiV1 "github.com/elizabeth-dev/ACME_Test/pkg/api/v1"
	"github.com/sirupsen/logrus"
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

const createUserTag = "CreateUser"

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
		if castErr, ok := err.(*errors.InvalidField); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": createUserTag,
					"cmd": cmd,
				},
			).WithError(castErr).Error("Invalid field")

			return nil, status.Error(codes.InvalidArgument, castErr.Error())
		}

		if castErr, ok := err.(*errors.MultipleInvalidFields); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": createUserTag,
					"cmd": cmd,
				},
			).WithError(castErr).Error("Invalid fields")

			return nil, status.Error(codes.InvalidArgument, castErr.Error())
		}

		logrus.WithFields(
			logrus.Fields{
				"tag": createUserTag,
				"cmd": cmd,
			},
		).WithError(err).Error("Unknown error while creating user")

		return nil, status.Error(codes.Internal, "Unknown error while creating user")
	}

	newUser, err := g.app.Queries.GetUserById.Handle(ctx, id)

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag": createUserTag,
				"id":  id,
			},
		).WithError(err).Error("Error retrieving new user")

		return nil, status.Error(codes.Unavailable, "The user was created but couldn't be retrieved")
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

const getUsersTag = "GetUsers"

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
		logrus.WithFields(
			logrus.Fields{
				"tag":   getUsersTag,
				"query": getUsersQuery,
			},
		).WithError(err).Error("Error retrieving users")

		return status.Error(codes.Internal, "Error retrieving users")
	}

	for i, currentUser := range users {
		if err := srv.Send(
			&apiV1.User{
				Id:        currentUser.Id,
				FirstName: currentUser.FirstName,
				LastName:  currentUser.LastName,
				Nickname:  currentUser.Nickname,
				Password:  currentUser.Password,
				Email:     currentUser.Email,
				Country:   currentUser.Country,
				CreatedAt: timestamppb.New(currentUser.CreatedAt),
				UpdatedAt: timestamppb.New(currentUser.UpdatedAt),
			},
		); err != nil {
			logrus.WithFields(
				logrus.Fields{
					"tag":   getUsersTag,
					"user":  currentUser,
					"index": i,
				},
			).WithError(err).Errorf("Error sending user")

			return status.Error(codes.Internal, "Error sending users")
		}
	}

	return nil
}

const updateUserTag = "UpdateUser"

func (g *GrpcServer) UpdateUser(ctx context.Context, request *apiV1.UpdateUserRequest) (*apiV1.User, error) {

	if request.GetId() == "" {
		logrus.WithFields(
			logrus.Fields{
				"tag":     updateUserTag,
				"request": request,
			},
		).Error("Error updating user: id is required")

		return nil, status.Error(codes.InvalidArgument, "Id is required")
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
		if castErr, ok := err.(*user.NotFoundError); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": updateUserTag,
					"cmd": cmd,
				},
			).WithError(castErr).Error("Attempted to update nonexistent user")

			return nil, status.Error(codes.NotFound, castErr.Error())
		}
		if castErr, ok := err.(*errors.InvalidField); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": updateUserTag,
					"cmd": cmd,
				},
			).WithError(castErr).Error("Invalid field")

			return nil, status.Error(codes.InvalidArgument, castErr.Error())
		}
		if castErr, ok := err.(*errors.MultipleInvalidFields); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": updateUserTag,
					"cmd": cmd,
				},
			).WithError(castErr).Error("Invalid fields")

			return nil, status.Error(codes.InvalidArgument, castErr.Error())
		}

		logrus.WithFields(
			logrus.Fields{
				"tag": updateUserTag,
				"cmd": cmd,
			},
		).WithError(err).Error("Unknown error while updating user")

		return nil, status.Error(codes.Internal, "Unknown error while updating user")
	}

	updatedUser, err := g.app.Queries.GetUserById.Handle(ctx, request.GetId())

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag": updateUserTag,
				"id":  request.GetId(),
			},
		).WithError(err).Error("Error retrieving updated user")

		return nil, status.Error(codes.Unavailable, "The user was updated but couldn't be retrieved")
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

const removeUserTag = "RemoveUser"

func (g *GrpcServer) RemoveUser(ctx context.Context, request *apiV1.RemoveUserRequest) (*emptypb.Empty, error) {
	if request.GetId() == "" {
		logrus.WithFields(
			logrus.Fields{
				"tag":     removeUserTag,
				"request": request,
			},
		).Error("Error removing user: id is required")

		return nil, status.Error(codes.InvalidArgument, "Id is required")
	}

	err := g.app.Commands.RemoveUser.Handle(ctx, request.GetId())

	if err != nil {
		if castErr, ok := err.(*user.NotFoundError); ok {
			logrus.WithFields(
				logrus.Fields{
					"tag": updateUserTag,
					"id":  request.GetId(),
				},
			).WithError(castErr).Error("Attempted to remove nonexistent user")

			return nil, status.Error(codes.NotFound, castErr.Error())
		}

		logrus.WithFields(
			logrus.Fields{
				"tag": removeUserTag,
				"id":  request.GetId(),
			},
		).WithError(err).Error("Error removing user")

		return nil, status.Error(codes.Internal, "Unknown error while removing user")
	}

	return &emptypb.Empty{}, nil
}
