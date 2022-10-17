package e2e

import (
	"context"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"testing"
)

func TestE2E(t *testing.T) {
	/* Setup */
	addr := os.Getenv("USER_URL")

	if addr == "" {
		t.Fatal("USER_URL env var is not set")
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Fatalf("could not close connection: %v", err)
		}
	}(conn)

	client := apiV1.NewUserServiceClient(conn)

	/* Tests */

	t.Run(
		"create users", func(t *testing.T) {
			testCreateUsersE2E(t, client)
		},
	)

	t.Run(
		"get created users", func(t *testing.T) {
			testGetCreatedUsersE2E(t, client)
		},
	)

	t.Run(
		"update users", func(t *testing.T) {
			testUpdateUsersE2E(t, client)
		},
	)

	t.Run(
		"get updated users", func(t *testing.T) {
			testGetUpdatedUsersE2E(t, client)
		},
	)

	t.Run(
		"remove users", func(t *testing.T) {
			testRemoveUsersE2E(t, client)
		},
	)

	t.Run(
		"get remaining users", func(t *testing.T) {
			testGetRemainingUsersE2E(t, client)
		},
	)
}

/* Utils */

func assertUserEquality(t *testing.T, expected *User, actual *apiV1.User) {
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.Nickname, actual.Nickname)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte(expected.Password)))
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Country, actual.Country)
}

func collectUsers(t *testing.T, client apiV1.UserService_GetUsersClient) (users []*apiV1.User) {
	for {
		user, err := client.Recv()

		if err == io.EOF {
			break
		}

		assert.NoError(t, err)
		assert.NotNil(t, user)

		users = append(users, user)
	}

	return users
}

func getSortedUsers(t *testing.T, client apiV1.UserServiceClient) (users []*apiV1.User) {
	prepareOut, _ := client.GetUsers(
		context.Background(), &apiV1.GetUsersRequest{
			Sort: []*apiV1.Sort{
				{
					Field:     "timestamp",
					Direction: apiV1.Sort_ASC,
				},
			},
		},
	)

	return collectUsers(t, prepareOut)
}
