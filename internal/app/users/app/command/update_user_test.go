package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	pkgErrors "github.com/elizabeth-dev/FACEIT_Test/internal/pkg/errors"
	"github.com/elizabeth-dev/FACEIT_Test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(t *testing.T){
		"initialize update user handler":                         testNewUpdateUserHandler,
		"initialize update user handler without repo":            testNewUpdateUserHandlerWithoutRepo,
		"handle update user command":                             testHandleUpdateUser,
		"handle update user command with user error":             testHandleUpdateUserWithUserError,
		"handle update user command with repo error on get user": testHandleUpdateUserWithRepoErrorOnGetUserById,
		"handle update user command with repo error on update":   testHandleUpdateUserWithRepoErrorOnUpdate,
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

func testNewUpdateUserHandler(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	newHandler := NewUpdateUserHandler(mockRepo)

	assert.NotNil(t, newHandler)
	assert.Equal(t, &UpdateUserHandler{mockRepo}, newHandler)
	assert.Same(t, mockRepo, newHandler.userRepo)
}

func testNewUpdateUserHandlerWithoutRepo(t *testing.T) {
	assert.PanicsWithValue(
		t, "[command/update_user] nil userRepo", func() {
			NewUpdateUserHandler(nil)
		},
	)
}

func testHandleUpdateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := UpdateUserHandler{mockRepo}

	ctx := context.Background()
	id := "123"
	previousUser := user.User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	updateCommand := UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updatedUser := previousUser
	_ = updatedUser.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	mockRepo.On("GetUserById", ctx, id).Return(&previousUser, nil)
	mockRepo.On(
		"UpdateUser", ctx, mock.MatchedBy(
			func(_user *user.User) bool {
				if (_user.Id() == updatedUser.Id()) &&
					(_user.FirstName() == updatedUser.FirstName()) &&
					(_user.LastName() == updatedUser.LastName()) &&
					(_user.Nickname() == updatedUser.Nickname()) &&
					(_user.Email() == updatedUser.Email()) &&
					(_user.Country() == updatedUser.Country()) {
					return true
				}

				return false
			},
		),
	).Return(nil)

	err := handler.Handle(ctx, updateCommand)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)

	assert.NoError(t, err)
}

func testHandleUpdateUserWithUserError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := UpdateUserHandler{mockRepo}

	ctx := context.Background()
	id := "123"
	previousUser := user.User1

	firstName := ""
	lastName := ""
	nickname := ""
	password := ""
	email := ""
	country := ""
	updateCommand := UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	mockRepo.On("GetUserById", ctx, id).Return(&previousUser, nil)

	err := handler.Handle(ctx, updateCommand)

	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)
	mockRepo.AssertNumberOfCalls(t, "UpdateUser", 0)
	mockRepo.AssertExpectations(t)

	assert.IsType(t, &pkgErrors.MultipleInvalidFields{}, err)

}

func testHandleUpdateUserWithRepoErrorOnGetUserById(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := UpdateUserHandler{mockRepo}

	ctx := context.Background()
	id := "123"
	previousUser := user.User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	updateCommand := UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updatedUser := previousUser
	_ = updatedUser.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	dbErr := errors.New("db is down")
	mockRepo.On("GetUserById", ctx, id).Return(nil, dbErr)

	err := handler.Handle(ctx, updateCommand)

	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)
	mockRepo.AssertNumberOfCalls(t, "UpdateUser", 0)
	mockRepo.AssertExpectations(t)

	assert.ErrorIs(t, err, dbErr)
}

func testHandleUpdateUserWithRepoErrorOnUpdate(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	handler := UpdateUserHandler{mockRepo}

	ctx := context.Background()
	id := "123"
	previousUser := user.User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	updateCommand := UpdateUser{
		Id:        id,
		FirstName: &firstName,
		LastName:  &lastName,
		Nickname:  &nickname,
		Password:  &password,
		Email:     &email,
		Country:   &country,
	}

	updatedUser := previousUser
	_ = updatedUser.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	dbErr := errors.New("db is down")
	mockRepo.On("GetUserById", ctx, id).Return(&previousUser, nil)
	mockRepo.On(
		"UpdateUser", ctx, mock.MatchedBy(
			func(_user *user.User) bool {
				if (_user.Id() == updatedUser.Id()) &&
					(_user.FirstName() == updatedUser.FirstName()) &&
					(_user.LastName() == updatedUser.LastName()) &&
					(_user.Nickname() == updatedUser.Nickname()) &&
					(_user.Email() == updatedUser.Email()) &&
					(_user.Country() == updatedUser.Country()) {
					return true
				}

				return false
			},
		),
	).Return(dbErr)

	err := handler.Handle(ctx, updateCommand)

	mockRepo.AssertNumberOfCalls(t, "GetUserById", 1)
	mockRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
	mockRepo.AssertExpectations(t)

	assert.ErrorIs(t, err, dbErr)
}
