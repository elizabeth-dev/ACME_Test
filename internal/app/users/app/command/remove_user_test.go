package command

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/ACME_Test/test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveUser(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize remove user handler":                  testNewRemoveUserHandler,
		"initialize remove user handler without repo":     testNewRemoveUserHandlerWithoutRepo,
		"handle remove user command":                      testHandleRemoveUser,
		"handle remove user command with error on get":    testHandleRemoveUserWithGetError,
		"handle remove user command with error on remove": testHandleRemoveUserWithRemoveError,
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

func testNewRemoveUserHandler(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	newHandler := NewRemoveUserHandler(mockRepo)

	assert.NotNil(t, newHandler)
	assert.Equal(t, &RemoveUserHandler{mockRepo}, newHandler)
	assert.Same(t, mockRepo, newHandler.userRepo)
}

func testNewRemoveUserHandlerWithoutRepo(t *testing.T) {
	assert.PanicsWithValue(
		t, "[command/remove_user] nil userRepo", func() {
			NewRemoveUserHandler(nil)
		},
	)
}

func testHandleRemoveUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := RemoveUserHandler{mockRepo}

	ctx := context.Background()

	removeId := user.User1.Id()

	mockRepo.On("GetUserById", ctx, removeId).Return(&user.User1, nil)
	mockRepo.On("RemoveUser", ctx, removeId).Return(nil)

	err := handler.Handle(ctx, removeId)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "RemoveUser", 1)
	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)

	assert.NoError(t, err)
}

func testHandleRemoveUserWithGetError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := RemoveUserHandler{mockRepo}

	ctx := context.Background()
	removeId := user.User1.Id()

	dbErr := errors.New("db is down")
	mockRepo.On("GetUserById", ctx, removeId).Return(nil, dbErr)

	err := handler.Handle(ctx, removeId)

	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)
	mockRepo.AssertNumberOfCalls(t, "RemoveUser", 0)
	mockRepo.AssertExpectations(t)

	assert.ErrorIs(t, err, dbErr)
}

func testHandleRemoveUserWithRemoveError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := RemoveUserHandler{mockRepo}

	ctx := context.Background()
	removeId := user.User1.Id()

	mockRepo.On("GetUserById", ctx, removeId).Return(&user.User1, nil)
	dbErr := errors.New("db is down")
	mockRepo.On("RemoveUser", ctx, removeId).Return(dbErr)

	err := handler.Handle(ctx, removeId)

	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)
	mockRepo.AssertNumberOfCalls(t, "RemoveUser", 1)
	mockRepo.AssertExpectations(t)

	assert.ErrorIs(t, err, dbErr)
}
