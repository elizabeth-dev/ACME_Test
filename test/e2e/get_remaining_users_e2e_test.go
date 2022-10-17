package e2e

import (
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func testGetRemainingUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	t.Run(
		"get remaining users", func(t *testing.T) {
			testRemainingUsers(t, client)
		},
	)
}

func testRemainingUsers(t *testing.T, client apiV1.UserServiceClient) {
	sortedUsers := getSortedUsers(t, client)

	require.Equal(t, 2, len(sortedUsers))
	assertUserEquality(t, &UpdatedUser0, sortedUsers[0])
	assertUserEquality(t, &User2, sortedUsers[1])
}
