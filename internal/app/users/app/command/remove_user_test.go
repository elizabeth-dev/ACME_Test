package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveUser(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize remove user handler":              testNewRemoveUserHandler,
		"initialize remove user handler without repo": testNewRemoveUserHandlerWithoutRepo,
		"handle remove user command":                  testHandleRemoveUser,
		"handle remove user command with repo error":  testHandleRemoveUserWithRepoError,
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

	removeId := "1234"

	mockRepo.On("RemoveUser", ctx, removeId).Return(nil)

	err := handler.Handle(ctx, removeId)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "RemoveUser", 1)

	assert.NoError(t, err)
}

func testHandleRemoveUserWithRepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := RemoveUserHandler{mockRepo}

	ctx := context.Background()
	removedId := "1234"

	mockRepo.On("RemoveUser", ctx, removedId).Return(errors.New("db is down"))

	err := handler.Handle(ctx, removedId)

	mockRepo.AssertNumberOfCalls(t, "RemoveUser", 1)
	mockRepo.AssertExpectations(t)

	assert.EqualError(t, err, "[command/remove_user] Error removing user 1234 from database: db is down")
}
