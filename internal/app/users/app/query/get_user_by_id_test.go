package query

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/ACME_Test/test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserById(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize get user by id handler":              testNewGetUserByIdHandler,
		"initialize get user by id handler without repo": testNewGetUserByIdHandlerWithoutRepo,
		"handle get user by id query":                    testHandleGetUserById,
		"handle get user by id query with repo error":    testHandleGetUserByIdWithRepoError,
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

func testNewGetUserByIdHandler(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	newHandler := NewGetUserByIdHandler(mockRepo)

	assert.NotNil(t, newHandler)
	assert.Equal(t, &GetUserByIdHandler{mockRepo}, newHandler)
	assert.Same(t, mockRepo, newHandler.userRepo)
}

func testNewGetUserByIdHandlerWithoutRepo(t *testing.T) {
	assert.PanicsWithValue(
		t, "[query/get_user_by_id] nil userRepo", func() {
			NewGetUserByIdHandler(nil)
		},
	)
}

func testHandleGetUserById(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := GetUserByIdHandler{mockRepo}

	ctx := context.Background()
	id := "1234"

	mockRepo.On("GetUserById", ctx, id).Return(&user.User1, nil)

	got, err := handler.Handle(ctx, id)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)

	assert.NoError(t, err)
	assert.Equal(
		t, &User{
			Id:        user.User1.Id(),
			FirstName: user.User1.FirstName(),
			LastName:  user.User1.LastName(),
			Nickname:  user.User1.Nickname(),
			Password:  user.User1.Password(),
			Email:     user.User1.Email(),
			Country:   user.User1.Country(),
			CreatedAt: user.User1.CreatedAt(),
			UpdatedAt: user.User1.UpdatedAt(),
		}, got,
	)
}

func testHandleGetUserByIdWithRepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := GetUserByIdHandler{mockRepo}

	ctx := context.Background()
	id := "1234"

	dbErr := errors.New("db is down")
	mockRepo.On("GetUserById", ctx, id).Return(nil, dbErr)

	got, err := handler.Handle(ctx, id)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}
