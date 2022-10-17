package e2e

import (
	"context"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func testGetCreatedUsersE2E(t *testing.T, client apiV1.UserServiceClient) {

	// Tests: Sort timestamp ASC
	t.Run(
		"get created users in ascending order", func(t *testing.T) {
			t.Parallel()

			testGetCreatedUsers(t, client)
		},
	)

	// Tests: Sort timestamp DESC
	t.Run(
		"get created users in descending order", func(t *testing.T) {
			t.Parallel()

			testGetCreatedUsersDescending(t, client)
		},
	)

	// Tests: Filter by string
	t.Run(
		"get created user by nickname", func(t *testing.T) {
			t.Parallel()

			testGetCreatedUserByNickname(t, client)
		},
	)

	// Tests: Limit + Filter by timestamp + Sort timestamp ASC
	t.Run(
		"get first created user only", func(t *testing.T) {
			t.Parallel()

			testGetFirstUserOnly(t, client)
		},
	)

	// Tests: Limit + Offset + Filter by timestamp + Sort timestamp ASC
	t.Run(
		"get third created user only", func(t *testing.T) {
			t.Parallel()

			testGetThirdUserOnly(t, client)
		},
	)

	t.Run(
		"check invalid user doesn't exist", func(t *testing.T) {
			t.Parallel()

			testGetInvalidUser(t, client)
		},
	)
}

func testGetCreatedUsers(t *testing.T, client apiV1.UserServiceClient) {
	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_ASC,
				},
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	require.Equal(t, 3, len(users))

	assertUserEquality(t, &User0, users[0])
	assert.Equal(t, users[0].UpdatedAt, users[0].CreatedAt)

	assertUserEquality(t, &User1, users[1])
	assert.Equal(t, users[1].UpdatedAt, users[1].CreatedAt)

	assertUserEquality(t, &User2, users[2])
	assert.Equal(t, users[2].UpdatedAt, users[2].CreatedAt)

	assert.True(t, users[0].CreatedAt.AsTime().Before(users[1].CreatedAt.AsTime()))
	assert.True(t, users[1].CreatedAt.AsTime().Before(users[2].CreatedAt.AsTime()))

}

func testGetCreatedUsersDescending(t *testing.T, client apiV1.UserServiceClient) {
	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_DESC,
				},
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	require.Equal(t, 3, len(users))

	assert.True(t, users[2].CreatedAt.AsTime().Before(users[1].CreatedAt.AsTime()))
	assert.True(t, users[1].CreatedAt.AsTime().Before(users[0].CreatedAt.AsTime()))

}

func testGetCreatedUserByNickname(t *testing.T, client apiV1.UserServiceClient) {
	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Filters: []*apiV1.Filter{
				{
					Field:    "nickname",
					Operator: apiV1.Filter_EQUALS,
					Value:    &apiV1.Filter_StringValue{StringValue: User1.Nickname},
				},
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	require.Equal(t, 1, len(users))
	assertUserEquality(t, &User1, users[0])
}

func testGetFirstUserOnly(t *testing.T, client apiV1.UserServiceClient) {
	prepareOut, _ := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_ASC,
				},
			},
		},
	)

	prepareUsers := collectUsers(t, prepareOut)

	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Filters: []*apiV1.Filter{
				{
					Field:    "created_at",
					Operator: apiV1.Filter_LESS_THAN,
					Value: &apiV1.Filter_TimestampValue{
						TimestampValue: &timestamppb.Timestamp{
							Seconds: prepareUsers[1].CreatedAt.Seconds,
							Nanos:   prepareUsers[1].CreatedAt.Nanos,
						},
					},
				},
			},
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_ASC,
				},
			},
			Pagination: &apiV1.Pagination{
				Limit: 1,
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	require.Equal(t, 1, len(users))
	assertUserEquality(t, &User0, users[0])
}

func testGetThirdUserOnly(t *testing.T, client apiV1.UserServiceClient) {
	prepareOut, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_ASC,
				},
			},
		},
	)

	prepareUsers := collectUsers(t, prepareOut)

	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Filters: []*apiV1.Filter{
				{
					Field:    "created_at",
					Operator: apiV1.Filter_GREATER_THAN_EQ,
					Value: &apiV1.Filter_TimestampValue{
						TimestampValue: &timestamppb.Timestamp{
							Seconds: prepareUsers[1].CreatedAt.Seconds,
							Nanos:   prepareUsers[1].CreatedAt.Nanos,
						},
					},
				},
			},
			Sort: []*apiV1.Sort{
				{
					Field:     "created_at",
					Direction: apiV1.Sort_ASC,
				},
			},
			Pagination: &apiV1.Pagination{
				Limit:  1,
				Offset: 1,
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	require.Equal(t, 1, len(users))
	assertUserEquality(t, &User2, users[0])
}

func testGetInvalidUser(t *testing.T, client apiV1.UserServiceClient) {
	out, err := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Filters: []*apiV1.Filter{
				{
					Field:    "nickname",
					Operator: apiV1.Filter_EQUALS,
					Value:    &apiV1.Filter_StringValue{StringValue: InvalidUser.Nickname},
				},
			},
		},
	)

	assert.NoError(t, err)

	users := collectUsers(t, out)

	assert.Equal(t, 0, len(users))
}
