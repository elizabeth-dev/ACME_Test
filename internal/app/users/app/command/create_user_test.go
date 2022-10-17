package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize create user handler":              testNewCreateUserHandler,
		"initialize create user handler without repo": testNewCreateUserHandlerWithoutRepo,
		"handle create user command":                  testHandleCreateUser,
		"handle create user command with user error":  testHandleCreateUserWithUserError,
		"handle create user command with repo error":  testHandleCreateUserWithRepoError,
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

func testNewCreateUserHandler(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	newHandler := NewCreateUserHandler(mockRepo)

	assert.NotNil(t, newHandler)
	assert.Equal(t, &CreateUserHandler{mockRepo}, newHandler)
	assert.Same(t, mockRepo, newHandler.userRepo)
}

func testNewCreateUserHandlerWithoutRepo(t *testing.T) {
	assert.PanicsWithValue(
		t, "[command/create_user] nil userRepo", func() {
			NewCreateUserHandler(nil)
		},
	)
}

func testHandleCreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := CreateUserHandler{mockRepo}

	ctx := context.Background()

	mockRepo.On("AddUser", ctx, mock.Anything).Return(nil)

	out, err := handler.Handle(
		ctx, CreateUser{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
		},
	)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "AddUser", 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, out)

}

func testHandleCreateUserWithUserError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := CreateUserHandler{mockRepo}

	ctx := context.Background()

	out, err := handler.Handle(
		ctx, CreateUser{
			FirstName: "",
			LastName:  "",
			Nickname:  "",
			Password:  "",
			Email:     "",
			Country:   "",
		},
	)

	mockRepo.AssertNumberOfCalls(t, "AddUser", 0)

	assert.ErrorContains(t, err, "[command/create_user] Error generating new user:")
	assert.Empty(t, out)

}

func testHandleCreateUserWithRepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := CreateUserHandler{mockRepo}

	ctx := context.Background()

	mockRepo.On("AddUser", ctx, mock.Anything).Return(errors.New("db is down"))

	out, err := handler.Handle(
		ctx, CreateUser{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "john-123",
			Password:  "password",
			Email:     "me@john.com",
			Country:   "US",
		},
	)

	mockRepo.AssertNumberOfCalls(t, "AddUser", 1)
	mockRepo.AssertExpectations(t)

	assert.EqualError(t, err, "[command/create_user] Error inserting new user into database: db is down")
	assert.Empty(t, out)
}
