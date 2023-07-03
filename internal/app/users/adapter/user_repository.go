package adapter

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/errors"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/helper/mongo_helper"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/mongo_utils"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

/*
UserModel holds the database representation of a user entity.

It may or may not align with the domain model, it should just hold the data
needed to generate the domain model.
*/
type UserModel struct {
	Id        string    `bson:"id"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Nickname  string    `bson:"nickname"`
	Password  string    `bson:"password"`
	Email     string    `bson:"email"`
	Country   string    `bson:"country"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserRepository struct {
	col mongo_helper.Collection
}

const UserRepoTag = "UserRepository"

func NewUserRepository(dbClient mongo_helper.Database) UserRepository {
	if dbClient == nil {
		log.Panicf("[%s] missing dbClient", UserRepoTag)
	}

	return UserRepository{col: dbClient.Collection("user")}
}

/*
AddUser inserts a whole user entity into the database.
*/
func (r *UserRepository) AddUser(ctx context.Context, newUser *user.User) error {
	logrus.WithFields(
		logrus.Fields{
			"tag":  UserRepoTag,
			"user": newUser,
		},
	).Debug("Adding user")
	userModel := r.marshalUser(newUser)

	if _, err := r.col.InsertOne(ctx, userModel); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":       UserRepoTag,
				"userModel": userModel,
			},
		).WithError(err).Error("Error inserting user")

		return &errors.Unknown{Tag: UserRepoTag, Cause: err}
	}

	return nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (*user.User, error) {
	logrus.WithFields(
		logrus.Fields{
			"tag":    UserRepoTag,
			"userId": userId,
		},
	).Debug("Getting user by id")
	var userModel UserModel

	if err := r.col.FindOne(ctx, bson.M{"id": userId}).Decode(&userModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &user.NotFoundError{Id: userId}
		}

		logrus.WithFields(
			logrus.Fields{
				"tag":    UserRepoTag,
				"userId": userId,
			},
		).WithError(err).Error("Error getting user by id")

		return nil, &errors.Unknown{Tag: UserRepoTag, Cause: err}
	}

	dbUser := user.UnmarshalUserFromDB(
		userModel.Id,
		userModel.FirstName,
		userModel.LastName,
		userModel.Nickname,
		userModel.Password,
		userModel.Email,
		userModel.Country,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return dbUser, nil
}

/*
GetUsers retrieves the user entities from the database.
*/
func (r *UserRepository) GetUsers(
	ctx context.Context, queryFilters []query_utils.Filter, sort []query_utils.Sort, pagination query_utils.Pagination,
) ([]*user.User, error) {
	logrus.WithFields(
		logrus.Fields{
			"tag":        UserRepoTag,
			"filters":    queryFilters,
			"sort":       sort,
			"pagination": pagination,
		},
	).Debug("Getting users")

	filter := mongo_utils.MapFilterToBson(queryFilters)
	opts := &options.FindOptions{
		Limit: &pagination.Limit,
		Skip:  &pagination.Offset,
	}

	if len(sort) > 0 {
		opts.SetSort(mongo_utils.MapSortToBson(sort))
	}

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":     UserRepoTag,
				"filter":  queryFilters,
				"options": opts,
			},
		).WithError(err).Error("Error getting users")

		return nil, &errors.Unknown{Tag: UserRepoTag, Cause: err}
	}

	var users []*user.User
	for cur.Next(ctx) {
		var userModel UserModel

		if err := cur.Decode(&userModel); err != nil {
			logrus.WithFields(
				logrus.Fields{
					"tag":     UserRepoTag,
					"filter":  queryFilters,
					"options": opts,
				},
			).WithError(err).Error("Error decoding user")

			return nil, &errors.Unknown{Tag: UserRepoTag, Cause: err}
		}

		dbUser := user.UnmarshalUserFromDB(
			userModel.Id,
			userModel.FirstName,
			userModel.LastName,
			userModel.Nickname,
			userModel.Password,
			userModel.Email,
			userModel.Country,
			userModel.CreatedAt,
			userModel.UpdatedAt,
		)

		users = append(users, dbUser)
	}

	return users, nil
}

/*
UpdateUser fully updates a user entity in the database.
*/
func (r *UserRepository) UpdateUser(ctx context.Context, userToUpdate *user.User) error {
	logrus.WithFields(
		logrus.Fields{
			"tag":  UserRepoTag,
			"user": userToUpdate,
		},
	).Debug("Updating user")
	userModel := r.marshalUser(userToUpdate)

	filter := bson.M{"id": userToUpdate.Id()}
	res, err := r.col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: userModel}})

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":       UserRepoTag,
				"filter":    filter,
				"userModel": userModel,
			},
		).WithError(err).Error("Error updating user")

		return &errors.Unknown{Tag: UserRepoTag, Cause: err}
	}

	if res.MatchedCount == 0 {
		return &user.NotFoundError{Id: userToUpdate.Id()}
	}

	return nil
}

/*
RemoveUser removes a user entity from the database given its id.
*/
func (r *UserRepository) RemoveUser(ctx context.Context, userId string) error {
	logrus.WithFields(
		logrus.Fields{
			"tag":    UserRepoTag,
			"userId": userId,
		},
	).Debug("Removing user")
	res, err := r.col.DeleteOne(ctx, bson.M{"id": userId})

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":    UserRepoTag,
				"userId": userId,
			},
		).WithError(err).Error("Error removing user")

		return &errors.Unknown{Tag: UserRepoTag, Cause: err}
	}

	if res.DeletedCount == 0 {
		return &user.NotFoundError{Id: userId}
	}

	return nil
}

/*
marshalUser converts a domain user entity into its database user model.
*/
func (r *UserRepository) marshalUser(user *user.User) *UserModel {
	return &UserModel{
		Id:        user.Id(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Nickname:  user.Nickname(),
		Password:  user.Password(),
		Email:     user.Email(),
		Country:   user.Country(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}
