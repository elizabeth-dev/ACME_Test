package e2e

import (
	"context"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func testCreateUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	t.Run(
		"create user 0", func(t *testing.T) {
			testCreateUser(t, client, User0)
		},
	)

	t.Run(
		"create user 1", func(t *testing.T) {
			testCreateUser(t, client, User1)
		},
	)

	t.Run(
		"create user 2", func(t *testing.T) {
			testCreateUser(t, client, User2)
		},
	)

	t.Run(
		"create invalid user", func(t *testing.T) {
			testCreateInvalidUser(t, client, InvalidUser)
		},
	)
}

func testCreateUser(t *testing.T, client apiV1.UserServiceClient, user User) {
	before := time.Now()

	out, err := client.CreateUser(
		context.Background(), &apiV1.CreateUserRequest{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
		},
	)

	after := time.Now()

	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, out.FirstName)
	assert.Equal(t, user.LastName, out.LastName)
	assert.Equal(t, user.Nickname, out.Nickname)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(out.Password), []byte(user.Password)))
	assert.Equal(t, user.Email, out.Email)
	assert.Equal(t, user.Country, out.Country)
	assert.True(t, before.Before(out.CreatedAt.AsTime()))
	assert.True(t, after.After(out.CreatedAt.AsTime()))
}

func testCreateInvalidUser(t *testing.T, client apiV1.UserServiceClient, user User) {
	out, err := client.CreateUser(
		context.Background(), &apiV1.CreateUserRequest{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
		},
	)

	assert.ErrorIs(
		t,
		err,
		status.Error(codes.Internal, "[command/create_user] Error generating new user: [User] Empty last name"),
	)
	assert.Nil(t, out)
}
