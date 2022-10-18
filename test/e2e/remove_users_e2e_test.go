package e2e

import (
	"context"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func testRemoveUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	t.Run(
		"remove user 1", func(t *testing.T) {
			t.Parallel()

			testRemoveUser1(t, client)
		},
	)

	t.Run(
		"remove invalid user", func(t *testing.T) {
			t.Parallel()

			testRemoveInvalidUser(t, client)
		},
	)

	t.Run(
		"remove nonexistent user", func(t *testing.T) {
			t.Parallel()

			testRemoveNonexistentUser(t, client)
		},
	)
}

func testRemoveUser1(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	out, err := client.RemoveUser(
		context.Background(), &apiV1.RemoveUserRequest{
			Id: sortedUsers[1].Id,
		},
	)

	assert.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, out)
}

func testRemoveInvalidUser(t *testing.T, client apiV1.UserServiceClient) {
	out, err := client.RemoveUser(
		context.Background(), &apiV1.RemoveUserRequest{
			Id: "",
		},
	)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "[RemoveUser] id is required"))
	assert.Nil(t, out)
}

func testRemoveNonexistentUser(t *testing.T, client apiV1.UserServiceClient) {
	id := "nonexistent"
	out, err := client.RemoveUser(
		context.Background(), &apiV1.RemoveUserRequest{
			Id: id,
		},
	)

	assert.ErrorIs(
		t,
		err,
		status.Error(
			codes.Internal,
			"[command/remove_user] Error retrieving user "+id+" from database: [UserRepository] User not found",
		),
	)
	assert.Nil(t, out)
}
