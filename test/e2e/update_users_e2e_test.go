package e2e

import (
	"context"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func testUpdateUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	t.Run(
		"update user 0", func(t *testing.T) {
			t.Parallel()

			testUpdateUser0(t, client)
		},
	)

	t.Run(
		"update user 1 with error", func(t *testing.T) {
			t.Parallel()

			testUpdateInvalidUser(t, client)
		},
	)

	t.Run(
		"update user 1 with several errors", func(t *testing.T) {
			t.Parallel()

			testUpdateUserWithSeveralErrors(t, client)
		},
	)

	t.Run(
		"update nonexistent user", func(t *testing.T) {
			t.Parallel()

			testUpdateNonexistentUser(t, client)
		},
	)
}

func testUpdateUser0(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	out, err := client.UpdateUser(
		context.Background(), &apiV1.UpdateUserRequest{
			Id:        sortedUsers[0].Id,
			FirstName: &UpdatedUser0.FirstName,
			LastName:  &UpdatedUser0.LastName,
			Nickname:  &UpdatedUser0.Nickname,
			Password:  &UpdatedUser0.Password,
			Email:     &UpdatedUser0.Email,
			Country:   &UpdatedUser0.Country,
		},
	)

	require.NoError(t, err)

	assertUserEquality(t, &UpdatedUser0, out)
}

func testUpdateInvalidUser(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	id := sortedUsers[1].Id
	out, err := client.UpdateUser(
		context.Background(), &apiV1.UpdateUserRequest{
			Id:        id,
			FirstName: &InvalidUpdatedUser0.FirstName,
			LastName:  &InvalidUpdatedUser0.LastName,
			Nickname:  &InvalidUpdatedUser0.Nickname,
			Password:  &InvalidUpdatedUser0.Password,
			Email:     &InvalidUpdatedUser0.Email,
			Country:   &InvalidUpdatedUser0.Country,
		},
	)

	assert.ErrorIs(
		t,
		err,
		status.Error(codes.Internal, "[command/update_user] Error updating user "+id+": [User] Empty password"),
	)
	assert.Nil(t, out)
}

func testUpdateUserWithSeveralErrors(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	id := sortedUsers[1].Id
	out, err := client.UpdateUser(
		context.Background(), &apiV1.UpdateUserRequest{
			Id:        id,
			FirstName: &InvalidUpdatedUser1.FirstName,
			LastName:  &InvalidUpdatedUser1.LastName,
			Nickname:  &InvalidUpdatedUser1.Nickname,
			Password:  &InvalidUpdatedUser1.Password,
			Email:     &InvalidUpdatedUser1.Email,
			Country:   &InvalidUpdatedUser1.Country,
		},
	)

	assert.ErrorIs(
		t,
		err,
		status.Error(codes.Internal, "[command/update_user] Error updating user "+id+": [User] Multiple errors"),
	)
	assert.Nil(t, out)
}

func testUpdateNonexistentUser(t *testing.T, client apiV1.UserServiceClient) {
	id := "nonexistent"
	out, err := client.UpdateUser(
		context.Background(), &apiV1.UpdateUserRequest{
			Id:        id,
			FirstName: &User2.FirstName,
			LastName:  &User2.LastName,
			Nickname:  &User2.Nickname,
			Password:  &User2.Password,
			Email:     &User2.Email,
			Country:   &User2.Country,
		},
	)

	assert.ErrorIs(
		t,
		err,
		status.Error(
			codes.Internal,
			"[command/update_user] Error getting user "+id+" from database: [UserRepository] User not found",
		),
	)
	assert.Nil(t, out)
}
