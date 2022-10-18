package query

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils/operators"
	"github.com/elizabeth-dev/FACEIT_Test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize get users handler":              testNewGetUsersHandler,
		"initialize get users handler without repo": testNewGetUsersHandlerWithoutRepo,
		"handle get users query":                    testHandleGetUsers,
		"handle get users query without parameters": testHandleGetUsersWithoutParameters,
		"handle get users query with repo error":    testHandleGetUsersWithRepoError,
	} {
		test := test
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()
				test(t)
			},
		)
	}
}

func testNewGetUsersHandler(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	newHandler := NewGetUsersHandler(mockRepo)

	assert.NotNil(t, newHandler)
	assert.Equal(t, &GetUsersHandler{mockRepo}, newHandler)
	assert.Same(t, mockRepo, newHandler.userRepo)
}

func testNewGetUsersHandlerWithoutRepo(t *testing.T) {
	assert.PanicsWithValue(
		t, "[query/get_users] nil userRepo", func() {
			NewGetUsersHandler(nil)
		},
	)
}

func testHandleGetUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := GetUsersHandler{mockRepo}

	ctx := context.Background()
	filters := []query_utils.Filter{
		{
			Field:    "id",
			Operator: operators.EQUALS,
			Value:    "123",
		},
	}
	sort := []query_utils.Sort{
		{
			Field:     "id",
			Direction: operators.ASC,
		},
	}
	pagination := query_utils.Pagination{
		Limit:  0,
		Offset: 0,
	}

	mockRepo.On("GetUsers", ctx, filters, sort, pagination).Return([]*user.User{&user.User1}, nil)

	out, err := handler.Handle(
		ctx, GetUsers{
			Filters:    filters,
			Sort:       sort,
			Pagination: pagination,
		},
	)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUsers", 1)

	assert.NoError(t, err)
	assert.Equal(
		t, []*User{
			{
				Id:        user.User1.Id(),
				FirstName: user.User1.FirstName(),
				LastName:  user.User1.LastName(),
				Nickname:  user.User1.Nickname(),
				Password:  user.User1.Password(),
				Email:     user.User1.Email(),
				Country:   user.User1.Country(),
				CreatedAt: user.User1.CreatedAt(),
				UpdatedAt: user.User1.UpdatedAt(),
			},
		}, out,
	)
}

func testHandleGetUsersWithoutParameters(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := GetUsersHandler{mockRepo}

	ctx := context.Background()
	var filters []query_utils.Filter
	var sort []query_utils.Sort
	pagination := query_utils.Pagination{
		Limit:  0,
		Offset: 0,
	}

	mockRepo.On("GetUsers", ctx, filters, sort, pagination).Return([]*user.User{&user.User1}, nil)

	out, err := handler.Handle(
		ctx, GetUsers{
			Filters:    filters,
			Sort:       sort,
			Pagination: pagination,
		},
	)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUsers", 1)

	assert.NoError(t, err)
	assert.Equal(
		t, []*User{
			{
				Id:        user.User1.Id(),
				FirstName: user.User1.FirstName(),
				LastName:  user.User1.LastName(),
				Nickname:  user.User1.Nickname(),
				Password:  user.User1.Password(),
				Email:     user.User1.Email(),
				Country:   user.User1.Country(),
				CreatedAt: user.User1.CreatedAt(),
				UpdatedAt: user.User1.UpdatedAt(),
			},
		}, out,
	)
}

func testHandleGetUsersWithRepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := GetUsersHandler{mockRepo}

	ctx := context.Background()
	filters := []query_utils.Filter{
		{
			Field:    "id",
			Operator: operators.EQUALS,
			Value:    "123",
		},
	}
	sort := []query_utils.Sort{
		{
			Field:     "id",
			Direction: operators.ASC,
		},
	}
	pagination := query_utils.Pagination{
		Limit:  0,
		Offset: 0,
	}

	dbErr := errors.New("db is down")
	mockRepo.On("GetUsers", ctx, filters, sort, pagination).Return(nil, dbErr)

	out, err := handler.Handle(
		ctx, GetUsers{
			Filters:    filters,
			Sort:       sort,
			Pagination: pagination,
		},
	)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUsers", 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, out)
}
