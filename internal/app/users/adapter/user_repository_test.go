package adapter

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/mongo_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils/operators"
	"github.com/elizabeth-dev/FACEIT_Test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var marshalledUser = UserModel{
	Id:        user.User1.Id(),
	FirstName: user.User1.FirstName(),
	LastName:  user.User1.LastName(),
	Nickname:  user.User1.Nickname(),
	Password:  user.User1.Password(),
	Email:     user.User1.Email(),
	Country:   user.User1.Country(),
	CreatedAt: user.User1.CreatedAt(),
	UpdatedAt: user.User1.UpdatedAt(),
}

func TestUserRepository(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"user repository": {
			"initialize user repository":                testNewUserRepository,
			"initialize user repository with no client": testNewUserRepositoryWithNoClient,
		},
		"add user": {
			"call add user":                  testAddUser,
			"call create user with db error": testAddUserWithDbError,
		},
		"get user by id": {
			"call get user by id":                     testGetUserById,
			"call get user by id with decode error":   testGetUserByIdWithDecodeError,
			"call get user by id with empty response": testGetUserByIdWithEmptyResponse,
		},
		"get users": {
			"call get users":                   testGetUsers,
			"call get users with no params":    testGetUsersWithNoParams,
			"call get users with db error":     testGetUsersWithDbError,
			"call get users with decode error": testGetUsersWithDecodeError,
		},
		"update user": {
			"call update user":               testUpdateUser,
			"call update user with db error": testUpdateUserWithDbError,
		},
		"remove user": {
			"call remove user":               testRemoveUser,
			"call remove user with db error": testRemoveUserWithDbError,
		},
		"marshal user": {
			"call marshal user": testMarshalUser,
		},
	} {
		testGroup := testGroup
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()

				for name, test := range testGroup {
					test := test
					t.Run(
						name, func(t *testing.T) {
							t.Parallel()

							test(t)
						},
					)
				}
			},
		)
	}

}

func testNewUserRepository(t *testing.T) {
	mockDb := new(mocks.Database)
	mockCol := new(mocks.Collection)

	mockDb.On("Collection", mock.Anything).Return(mockCol)

	out := NewUserRepository(mockDb)

	assert.NotNil(t, out)
	assert.Equal(t, mockCol, out.col)
}

func testNewUserRepositoryWithNoClient(t *testing.T) {
	assert.PanicsWithValue(
		t, "[UserRepository] missing dbClient", func() {
			NewUserRepository(nil)
		},
	)
}

func testAddUser(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	newUser := &user.User1

	mockCollection.On("InsertOne", ctx, &marshalledUser).Return(nil, nil)

	err := repo.AddUser(ctx, newUser)

	mockCollection.AssertNumberOfCalls(t, "InsertOne", 1)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
}

func testAddUserWithDbError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	newUser := &user.User1

	mockCollection.On("InsertOne", ctx, &marshalledUser).Return(nil, errors.New("db error"))

	err := repo.AddUser(ctx, newUser)

	mockCollection.AssertNumberOfCalls(t, "InsertOne", 1)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error creating user: db error")
}

func testGetUserById(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockSingleResult := new(mocks.SingleResult)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	id := "1234"

	mockCollection.On("FindOne", ctx, bson.M{"id": id}).Return(mockSingleResult, nil)
	mockSingleResult.On("Decode", &UserModel{}).Run(
		func(args mock.Arguments) {
			*args.Get(0).(*UserModel) = marshalledUser
		},
	).Return(nil)

	out, err := repo.GetUserById(ctx, id)

	mockSingleResult.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "FindOne", 1)
	mockSingleResult.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, &user.User1, out)
}

func testGetUserByIdWithDecodeError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockSingleResult := new(mocks.SingleResult)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	id := "1234"

	mockCollection.On("FindOne", ctx, bson.M{"id": id}).Return(mockSingleResult, nil)
	mockSingleResult.On("Decode", &UserModel{}).Return(errors.New("decode error"))

	out, err := repo.GetUserById(ctx, id)

	mockSingleResult.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "FindOne", 1)
	mockSingleResult.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error retrieving user: decode error")
	assert.Nil(t, out)
}

func testGetUserByIdWithEmptyResponse(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockSingleResult := new(mocks.SingleResult)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	id := "1234"

	mockCollection.On("FindOne", ctx, bson.M{"id": id}).Return(mockSingleResult, nil)
	mockSingleResult.On("Decode", &UserModel{}).Return(mongo.ErrNoDocuments)

	out, err := repo.GetUserById(ctx, id)

	mockSingleResult.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "FindOne", 1)
	mockSingleResult.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] User not found")
	assert.Nil(t, out)
}

func testGetUsers(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockCursor := new(mocks.Cursor)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()
	filters := []query_utils.Filter{
		{
			Field:    "id",
			Operator: operators.EQUALS,
			Value:    "1234",
		},
	}
	sort := []query_utils.Sort{
		{
			Field:     "id",
			Direction: operators.ASC,
		},
	}
	pagination := query_utils.Pagination{
		Limit:  10,
		Offset: 0,
	}

	findOptions := options.FindOptions{
		Sort:  mongo_utils.MapSortToBson(sort),
		Limit: &pagination.Limit,
		Skip:  &pagination.Offset,
	}

	mockCollection.On("Find", ctx, mongo_utils.MapFilterToBson(filters), &findOptions).Return(mockCursor, nil)
	mockCursor.On("Next", ctx).Return(true).Once()
	mockCursor.On("Decode", &UserModel{}).Run(
		func(args mock.Arguments) {
			*args.Get(0).(*UserModel) = marshalledUser
		},
	).Return(nil).Once()
	mockCursor.On("Next", ctx).Return(false)

	out, err := repo.GetUsers(ctx, filters, sort, pagination)

	mockCursor.AssertNumberOfCalls(t, "Next", 2)
	mockCursor.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "Find", 1)
	mockCursor.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, []*user.User{&user.User1}, out)
}

func testGetUsersWithNoParams(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockCursor := new(mocks.Cursor)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	zero := int64(0)

	findOptions := options.FindOptions{
		Limit: &zero,
		Skip:  &zero,
	}

	mockCollection.On("Find", ctx, bson.M{}, &findOptions).Return(mockCursor, nil)
	mockCursor.On("Next", ctx).Return(true).Once()
	mockCursor.On("Decode", &UserModel{}).Run(
		func(args mock.Arguments) {
			*args.Get(0).(*UserModel) = marshalledUser
		},
	).Return(nil).Once()
	mockCursor.On("Next", ctx).Return(false)

	out, err := repo.GetUsers(ctx, nil, nil, query_utils.Pagination{})

	mockCursor.AssertNumberOfCalls(t, "Next", 2)
	mockCursor.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "Find", 1)
	mockCursor.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, []*user.User{&user.User1}, out)
}

func testGetUsersWithDbError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockCursor := new(mocks.Cursor)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	zero := int64(0)

	findOptions := options.FindOptions{
		Limit: &zero,
		Skip:  &zero,
	}

	mockCollection.On("Find", ctx, bson.M{}, &findOptions).Return(nil, errors.New("db error"))

	out, err := repo.GetUsers(ctx, nil, nil, query_utils.Pagination{})

	mockCursor.AssertNumberOfCalls(t, "Next", 0)
	mockCursor.AssertNumberOfCalls(t, "Decode", 0)
	mockCollection.AssertNumberOfCalls(t, "Find", 1)
	mockCursor.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error retrieving users: db error")
	assert.Nil(t, out)
}

func testGetUsersWithDecodeError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	mockCursor := new(mocks.Cursor)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	zero := int64(0)

	findOptions := options.FindOptions{
		Limit: &zero,
		Skip:  &zero,
	}

	mockCollection.On("Find", ctx, bson.M{}, &findOptions).Return(mockCursor, nil)
	mockCursor.On("Next", ctx).Return(true).Once()
	mockCursor.On("Decode", &UserModel{}).Return(errors.New("decode error"))

	out, err := repo.GetUsers(ctx, nil, nil, query_utils.Pagination{})

	mockCursor.AssertNumberOfCalls(t, "Next", 1)
	mockCursor.AssertNumberOfCalls(t, "Decode", 1)
	mockCollection.AssertNumberOfCalls(t, "Find", 1)
	mockCursor.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error decoding user: decode error")
	assert.Nil(t, out)
}

func testUpdateUser(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	mockCollection.On(
		"UpdateOne", ctx, bson.M{"id": user.User1.Id()}, bson.D{{Key: "$set", Value: &marshalledUser}},
	).Return(nil, nil)

	err := repo.UpdateUser(ctx, &user.User1)

	mockCollection.AssertNumberOfCalls(t, "UpdateOne", 1)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
}

func testUpdateUserWithDbError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	mockCollection.On(
		"UpdateOne", ctx, bson.M{"id": user.User1.Id()}, bson.D{{Key: "$set", Value: &marshalledUser}},
	).Return(
		nil, errors.New("db error"),
	)

	err := repo.UpdateUser(ctx, &user.User1)

	mockCollection.AssertNumberOfCalls(t, "UpdateOne", 1)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error updating user "+user.User1.Id()+": db error")
}

func testRemoveUser(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	mockCollection.On("DeleteOne", ctx, bson.M{"id": user.User1.Id()}).Return(int64(0), nil)

	err := repo.RemoveUser(ctx, user.User1.Id())

	mockCollection.AssertNumberOfCalls(t, "DeleteOne", 1)
	mockCollection.AssertExpectations(t)

	assert.NoError(t, err)
}

func testRemoveUserWithDbError(t *testing.T) {
	mockCollection := new(mocks.Collection)
	repo := UserRepository{col: mockCollection}

	ctx := context.Background()

	mockCollection.On("DeleteOne", ctx, bson.M{"id": user.User1.Id()}).Return(int64(0), errors.New("db error"))

	err := repo.RemoveUser(ctx, user.User1.Id())

	mockCollection.AssertNumberOfCalls(t, "DeleteOne", 1)
	mockCollection.AssertExpectations(t)

	assert.EqualError(t, err, "[UserRepository] Error deleting user "+user.User1.Id()+": db error")
}

func testMarshalUser(t *testing.T) {
	repo := UserRepository{}

	out := repo.marshalUser(&user.User1)

	assert.Equal(t, &marshalledUser, out)
}
