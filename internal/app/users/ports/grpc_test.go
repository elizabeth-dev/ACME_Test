package ports

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/command"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/query"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	errors2 "github.com/elizabeth-dev/ACME_Test/internal/pkg/errors"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils/operators"
	apiV1 "github.com/elizabeth-dev/ACME_Test/pkg/api/v1"
	"github.com/elizabeth-dev/ACME_Test/test/mocks"
	handler_mocks2 "github.com/elizabeth-dev/ACME_Test/test/mocks/handler_mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"gRPC server": {
			"initialize gRPC server": testNewGrpcServer,
		},
		"create user": {
			"call create user":                                    testCreateUser,
			"call create user with invalid field error":           testCreateUserWithInvalidFieldError,
			"call create user with multiple invalid fields error": testCreateUserWithMultipleInvalidFieldsError,
			"call create user with create error":                  testCreateUserWithCreateError,
			"call create user with get error":                     testCreateUserWithGetError,
		},
		"get users": {
			"call get users":                    testGetUsers,
			"call get users with no parameters": testGetUsersWithNoParams,
			"call get users with get error":     testGetUsersWithGetError,
			"call get users with send error":    testGetUsersWithSendError,
		},
		"update user": {
			"call update user":                                   testUpdateUser,
			"call update user with no id":                        testUpdateUserWithoutId,
			"call update user with not found error":              testUpdateUserWithNotFoundError,
			"call update user with invalid field error":          testUpdateUserWithInvalidFieldError,
			"call update user with multiple invalidFields error": testUpdateUserWithMultipleInvalidFieldsError,
			"call update user with update error":                 testUpdateUserWithUpdateError,
			"call update user with get error":                    testUpdateUserWithGetError,
		},
		"remove user": {
			"call remove user":                      testRemoveUser,
			"call remove user with no id":           testRemoveUserWithoutId,
			"call remove user with not found error": testRemoveUserWithNotFoundError,
			"call remove user with remove error":    testRemoveUserWithRemoveError,
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

func testNewGrpcServer(t *testing.T) {
	application := app.Application{}

	out := NewGrpcServer(application)

	assert.NotNil(t, out)
	assert.Equal(t, GrpcServer{application}, out)
}

func testCreateUser(t *testing.T) {
	mockCreateUser := new(handler_mocks2.ICreateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{CreateUser: mockCreateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	createUserCmd := command.CreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	now := time.Now()
	getUserResult := query.User{
		Id:        id,
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockCreateUser.On("Handle", ctx, createUserCmd).Return(id, nil)
	mockGetUserById.On("Handle", ctx, id).Return(&getUserResult, nil)

	out, err := server.CreateUser(ctx, &request)

	mockCreateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 1)
	mockCreateUser.AssertExpectations(t)
	mockGetUserById.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(
		t, &apiV1.User{
			Id:        getUserResult.Id,
			FirstName: getUserResult.FirstName,
			LastName:  getUserResult.LastName,
			Nickname:  getUserResult.Nickname,
			Password:  getUserResult.Password,
			Email:     getUserResult.Email,
			Country:   getUserResult.Country,
			CreatedAt: timestamppb.New(getUserResult.CreatedAt),
			UpdatedAt: timestamppb.New(getUserResult.UpdatedAt),
		}, out,
	)
}

func testCreateUserWithInvalidFieldError(t *testing.T) {
	mockCreateUser := new(handler_mocks2.ICreateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{CreateUser: mockCreateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	ctx := context.Background()
	request := apiV1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	createUserCmd := command.CreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	invalidErr := errors2.InvalidField{
		Domain: "",
		Field:  "",
		Value:  nil,
	}
	mockCreateUser.On("Handle", ctx, createUserCmd).Return("", &invalidErr)

	out, err := server.CreateUser(ctx, &request)

	mockCreateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockCreateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, invalidErr.Error()))
	assert.Nil(t, out)
}

func testCreateUserWithMultipleInvalidFieldsError(t *testing.T) {
	mockCreateUser := new(handler_mocks2.ICreateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{CreateUser: mockCreateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	ctx := context.Background()
	request := apiV1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	createUserCmd := command.CreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	invalidErr := errors2.MultipleInvalidFields{
		Errors: []error{
			&errors2.InvalidField{
				Domain: "",
				Field:  "",
				Value:  nil,
			},
		},
	}
	mockCreateUser.On("Handle", ctx, createUserCmd).Return(
		"", &invalidErr,
	)

	out, err := server.CreateUser(ctx, &request)

	mockCreateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockCreateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, invalidErr.Error()))
	assert.Nil(t, out)
}

func testCreateUserWithCreateError(t *testing.T) {
	mockCreateUser := new(handler_mocks2.ICreateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{CreateUser: mockCreateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	ctx := context.Background()
	request := apiV1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	createUserCmd := command.CreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	mockCreateUser.On("Handle", ctx, createUserCmd).Return("", errors.New("unknown error"))

	out, err := server.CreateUser(ctx, &request)

	mockCreateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockCreateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "Unknown error while creating user"))
	assert.Nil(t, out)
}

func testCreateUserWithGetError(t *testing.T) {
	mockCreateUser := new(handler_mocks2.ICreateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{CreateUser: mockCreateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	createUserCmd := command.CreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "john-123",
		Password:  "password",
		Email:     "me@john.com",
		Country:   "US",
	}

	mockCreateUser.On("Handle", ctx, createUserCmd).Return(id, nil)
	mockGetUserById.On("Handle", ctx, id).Return(nil, errors.New("unknown error"))

	out, err := server.CreateUser(ctx, &request)

	mockCreateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 1)
	mockCreateUser.AssertExpectations(t)
	mockGetUserById.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Unavailable, "The user was created but couldn't be retrieved"))
	assert.Nil(t, out)
}

func testGetUsers(t *testing.T) {
	mockGetUsersHandler := new(handler_mocks2.IGetUsersHandler)
	mockGetUsersSrv := new(mocks.UserService_GetUsersServer)
	application := app.Application{
		Queries: app.Queries{GetUsers: mockGetUsersHandler},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.GetUsersRequest{
		Filters: []*apiV1.Filter{
			{
				Field:    "id",
				Operator: apiV1.Filter_EQUALS,
				Value:    &apiV1.Filter_StringValue{StringValue: ""},
			},
		},
		Sort: []*apiV1.Sort{
			{
				Field:     "id",
				Direction: apiV1.Sort_ASC,
			},
		},
		Pagination: &apiV1.Pagination{
			Limit:  0,
			Offset: 0,
		},
	}

	getUsersQuery := query.GetUsers{
		Filters: []query_utils.Filter{
			{
				Field:    "id",
				Operator: operators.EQUALS,
				Value:    "",
			},
		},
		Sort: []query_utils.Sort{
			{
				Field:     "id",
				Direction: operators.ASC,
			},
		},
		Pagination: query_utils.Pagination{},
	}

	now := time.Now()
	getUsersResult := []*query.User{
		{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockGetUsersSrv.On("Context").Return(ctx)
	mockGetUsersSrv.On(
		"Send", &apiV1.User{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	).Return(nil)
	mockGetUsersHandler.On("Handle", ctx, getUsersQuery).Return(getUsersResult, nil)

	err := server.GetUsers(&request, mockGetUsersSrv)

	mockGetUsersSrv.AssertNumberOfCalls(t, "Context", 1)
	mockGetUsersSrv.AssertNumberOfCalls(t, "Send", len(getUsersResult))
	mockGetUsersHandler.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUsersHandler.AssertExpectations(t)
	mockGetUsersSrv.AssertExpectations(t)

	assert.NoError(t, err)
}

func testGetUsersWithNoParams(t *testing.T) {
	mockGetUsersHandler := new(handler_mocks2.IGetUsersHandler)
	mockGetUsersSrv := new(mocks.UserService_GetUsersServer)
	application := app.Application{
		Queries: app.Queries{GetUsers: mockGetUsersHandler},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.GetUsersRequest{
		Filters:    nil,
		Sort:       nil,
		Pagination: nil,
	}

	getUsersQuery := query.GetUsers{
		Filters: nil,
		Sort:    nil,
		Pagination: query_utils.Pagination{
			Limit:  0,
			Offset: 0,
		},
	}

	now := time.Now()
	getUsersResult := []*query.User{
		{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockGetUsersSrv.On("Context").Return(ctx)
	mockGetUsersSrv.On(
		"Send", &apiV1.User{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	).Return(nil)
	mockGetUsersHandler.On("Handle", ctx, getUsersQuery).Return(getUsersResult, nil)

	err := server.GetUsers(&request, mockGetUsersSrv)

	mockGetUsersSrv.AssertNumberOfCalls(t, "Context", 1)
	mockGetUsersSrv.AssertNumberOfCalls(t, "Send", len(getUsersResult))
	mockGetUsersHandler.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUsersHandler.AssertExpectations(t)
	mockGetUsersSrv.AssertExpectations(t)

	assert.NoError(t, err)
}

func testGetUsersWithGetError(t *testing.T) {
	mockGetUsersHandler := new(handler_mocks2.IGetUsersHandler)
	mockGetUsersSrv := new(mocks.UserService_GetUsersServer)
	application := app.Application{
		Queries: app.Queries{GetUsers: mockGetUsersHandler},
	}
	server := GrpcServer{app: application}

	ctx := context.Background()
	request := apiV1.GetUsersRequest{
		Filters: []*apiV1.Filter{
			{
				Field:    "id",
				Operator: apiV1.Filter_EQUALS,
				Value:    &apiV1.Filter_StringValue{StringValue: ""},
			},
		},
		Sort: []*apiV1.Sort{
			{
				Field:     "id",
				Direction: apiV1.Sort_ASC,
			},
		},
		Pagination: &apiV1.Pagination{
			Limit:  0,
			Offset: 0,
		},
	}

	getUsersQuery := query.GetUsers{
		Filters: []query_utils.Filter{
			{
				Field:    "id",
				Operator: operators.EQUALS,
				Value:    "",
			},
		},
		Sort: []query_utils.Sort{
			{
				Field:     "id",
				Direction: operators.ASC,
			},
		},
		Pagination: query_utils.Pagination{},
	}

	mockGetUsersSrv.On("Context").Return(ctx)
	mockGetUsersHandler.On("Handle", ctx, getUsersQuery).Return(nil, errors.New("unknown error"))

	err := server.GetUsers(&request, mockGetUsersSrv)

	mockGetUsersSrv.AssertNumberOfCalls(t, "Context", 1)
	mockGetUsersHandler.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUsersHandler.AssertExpectations(t)
	mockGetUsersSrv.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "Error retrieving users"))
}

func testGetUsersWithSendError(t *testing.T) {
	mockGetUsersHandler := new(handler_mocks2.IGetUsersHandler)
	mockGetUsersSrv := new(mocks.UserService_GetUsersServer)
	application := app.Application{
		Queries: app.Queries{GetUsers: mockGetUsersHandler},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.GetUsersRequest{
		Filters: []*apiV1.Filter{
			{
				Field:    "id",
				Operator: apiV1.Filter_EQUALS,
				Value:    &apiV1.Filter_StringValue{StringValue: ""},
			},
		},
		Sort: []*apiV1.Sort{
			{
				Field:     "id",
				Direction: apiV1.Sort_ASC,
			},
		},
		Pagination: &apiV1.Pagination{
			Limit:  0,
			Offset: 0,
		},
	}

	getUsersQuery := query.GetUsers{
		Filters: []query_utils.Filter{
			{
				Field:    "id",
				Operator: operators.EQUALS,
				Value:    "",
			},
		},
		Sort: []query_utils.Sort{
			{
				Field:     "id",
				Direction: operators.ASC,
			},
		},
		Pagination: query_utils.Pagination{},
	}

	now := time.Now()
	getUsersResult := []*query.User{
		{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockGetUsersSrv.On("Context").Return(ctx)
	mockGetUsersSrv.On(
		"Send", &apiV1.User{
			Id:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	).Return(errors.New("unknown error"))
	mockGetUsersHandler.On("Handle", ctx, getUsersQuery).Return(getUsersResult, nil)

	err := server.GetUsers(&request, mockGetUsersSrv)

	mockGetUsersSrv.AssertNumberOfCalls(t, "Context", 1)
	mockGetUsersSrv.AssertNumberOfCalls(t, "Send", 1)
	mockGetUsersHandler.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUsersHandler.AssertExpectations(t)
	mockGetUsersSrv.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "Error sending users"))
}

func testUpdateUser(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	now := time.Now()
	getUserResult := query.User{
		Id:        id,
		FirstName: "updated",
		LastName:  "updated",
		Nickname:  "updated",
		Password:  "updated",
		Email:     "updated",
		Country:   "updates",
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(nil)
	mockGetUserById.On("Handle", ctx, id).Return(&getUserResult, nil)

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 1)
	mockUpdateUser.AssertExpectations(t)
	mockGetUserById.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(
		t, &apiV1.User{
			Id:        getUserResult.Id,
			FirstName: getUserResult.FirstName,
			LastName:  getUserResult.LastName,
			Nickname:  getUserResult.Nickname,
			Password:  getUserResult.Password,
			Email:     getUserResult.Email,
			Country:   getUserResult.Country,
			CreatedAt: timestamppb.New(getUserResult.CreatedAt),
			UpdatedAt: timestamppb.New(getUserResult.UpdatedAt),
		}, out,
	)
}

func testUpdateUserWithoutId(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := ""
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 0)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "Id is required"))
	assert.Nil(t, out)
}

func testUpdateUserWithNotFoundError(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	notFoundErr := user.NotFoundError{Id: id}
	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(&notFoundErr)

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockUpdateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.NotFound, notFoundErr.Error()))
	assert.Nil(t, out)
}

func testUpdateUserWithInvalidFieldError(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	invalidErr := errors2.InvalidField{
		Domain: "",
		Field:  "",
		Value:  nil,
	}
	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(&invalidErr)

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockUpdateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, invalidErr.Error()))
	assert.Nil(t, out)
}

func testUpdateUserWithMultipleInvalidFieldsError(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	invalidErr := errors2.MultipleInvalidFields{
		Errors: []error{
			&errors2.InvalidField{
				Domain: "",
				Field:  "",
				Value:  nil,
			},
		},
	}
	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(
		&invalidErr,
	)

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockUpdateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, invalidErr.Error()))
	assert.Nil(t, out)
}

func testUpdateUserWithUpdateError(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(errors.New("unknown error"))

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 0)
	mockUpdateUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "Unknown error while updating user"))
	assert.Nil(t, out)
}

func testUpdateUserWithGetError(t *testing.T) {
	mockUpdateUser := new(handler_mocks2.IUpdateUserHandler)
	mockGetUserById := new(handler_mocks2.IGetUserByIdHandler)
	application := app.Application{
		Commands: app.Commands{UpdateUser: mockUpdateUser},
		Queries:  app.Queries{GetUserById: mockGetUserById},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	request := apiV1.UpdateUserRequest{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updateUserCmd := command.UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	mockUpdateUser.On("Handle", ctx, updateUserCmd).Return(nil)
	mockGetUserById.On("Handle", ctx, id).Return(nil, errors.New("unknown error"))

	out, err := server.UpdateUser(ctx, &request)

	mockUpdateUser.AssertNumberOfCalls(t, "Handle", 1)
	mockGetUserById.AssertNumberOfCalls(t, "Handle", 1)
	mockUpdateUser.AssertExpectations(t)
	mockGetUserById.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Unavailable, "The user was updated but couldn't be retrieved"))
	assert.Nil(t, out)
}

func testRemoveUser(t *testing.T) {
	mockRemoveUser := new(handler_mocks2.IRemoveUserHandler)
	application := app.Application{
		Commands: app.Commands{RemoveUser: mockRemoveUser},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.RemoveUserRequest{
		Id: id,
	}

	mockRemoveUser.On("Handle", ctx, id).Return(nil)

	out, err := server.RemoveUser(ctx, &request)

	mockRemoveUser.AssertNumberOfCalls(t, "Handle", 1)
	mockRemoveUser.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(
		t, &emptypb.Empty{}, out,
	)
}

func testRemoveUserWithoutId(t *testing.T) {
	mockRemoveUser := new(handler_mocks2.IRemoveUserHandler)
	application := app.Application{
		Commands: app.Commands{RemoveUser: mockRemoveUser},
	}
	server := GrpcServer{app: application}

	id := ""
	ctx := context.Background()
	request := apiV1.RemoveUserRequest{Id: id}

	out, err := server.RemoveUser(ctx, &request)

	mockRemoveUser.AssertNumberOfCalls(t, "Handle", 0)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "Id is required"))
	assert.Nil(t, out)
}

func testRemoveUserWithNotFoundError(t *testing.T) {
	mockRemoveUser := new(handler_mocks2.IRemoveUserHandler)
	application := app.Application{
		Commands: app.Commands{RemoveUser: mockRemoveUser},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.RemoveUserRequest{Id: id}

	notFoundErr := user.NotFoundError{Id: id}
	mockRemoveUser.On("Handle", ctx, id).Return(&notFoundErr)

	out, err := server.RemoveUser(ctx, &request)

	mockRemoveUser.AssertNumberOfCalls(t, "Handle", 1)
	mockRemoveUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.NotFound, notFoundErr.Error()))
	assert.Nil(t, out)
}

func testRemoveUserWithRemoveError(t *testing.T) {
	mockRemoveUser := new(handler_mocks2.IRemoveUserHandler)
	application := app.Application{
		Commands: app.Commands{RemoveUser: mockRemoveUser},
	}
	server := GrpcServer{app: application}

	id := "1234"
	ctx := context.Background()
	request := apiV1.RemoveUserRequest{Id: id}

	mockRemoveUser.On("Handle", ctx, id).Return(errors.New("unknown error"))

	out, err := server.RemoveUser(ctx, &request)

	mockRemoveUser.AssertNumberOfCalls(t, "Handle", 1)
	mockRemoveUser.AssertExpectations(t)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "Unknown error while removing user"))
	assert.Nil(t, out)
}
