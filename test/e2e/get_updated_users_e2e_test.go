package e2e

import (
	apiV1 "github.com/elizabeth-dev/ACME_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testGetUpdatedUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	t.Run(
		"get updated user 0", func(t *testing.T) {
			t.Parallel()

			testUpdatedUser0(t, client)
		},
	)

	t.Run(
		"get non-updated user 1", func(t *testing.T) {
			t.Parallel()

			testUntouchedUser1(t, client)
		},
	)

	t.Run(
		"get non-updated user 2", func(t *testing.T) {
			t.Parallel()

			testUntouchedUser2(t, client)
		},
	)
}

func testUpdatedUser0(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	assertUserEquality(t, &UpdatedUser0, sortedUsers[0])
	assert.True(t, sortedUsers[0].UpdatedAt.AsTime().After(sortedUsers[0].CreatedAt.AsTime()))
}

func testUntouchedUser1(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	assertUserEquality(t, &User1, sortedUsers[1])
	assert.Equal(t, sortedUsers[1].UpdatedAt.AsTime(), sortedUsers[1].CreatedAt.AsTime())
}

func testUntouchedUser2(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	assertUserEquality(t, &User2, sortedUsers[2])
	assert.Equal(t, sortedUsers[2].UpdatedAt.AsTime(), sortedUsers[2].CreatedAt.AsTime())
}
