package adapter

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	col *mongo.Collection
}

func NewUserRepository(dbClient *mongo.Database) UserRepository {
	if dbClient == nil {
		panic("[UserRepository] missing dbClient")
	}

	return UserRepository{col: dbClient.Collection("user")}
}

/*
AddUser inserts a whole user entity into the database.
*/
func (r *UserRepository) AddUser(ctx context.Context, newUser *user.User) error {
	userModel := r.marshalUser(newUser)

	if _, err := r.col.InsertOne(ctx, userModel); err != nil {
		return err
	}

	return nil
}

/*
GetUsers retrieves the user entities from the database.
*/
func (r *UserRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	var users []*user.User

	cur, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "[UserRepository] Error retrieving users")
	}

	for cur.Next(ctx) {
		var userModel UserModel

		if err := cur.Decode(&userModel); err != nil {
			return nil, errors.Wrap(err, "[UserRepository] Error decoding user")
		}

		u := user.UnmarshalUserFromDB(
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

		users = append(users, u)
	}

	return users, nil
}

/*
UpdateUser fully updates a user entity in the database.
*/
func (r *UserRepository) UpdateUser(ctx context.Context, user *user.User) error {
	userModel := r.marshalUser(user)

	_, err := r.col.UpdateOne(ctx, bson.M{"id": user.Id()}, bson.D{{Key: "$set", Value: userModel}})

	if err != nil {
		return errors.Wrap(err, "[UserRepository] Error updating user "+user.Id())
	}

	return nil
}

/*
RemoveUser removes a user entity from the database given its id.
*/
func (r *UserRepository) RemoveUser(ctx context.Context, userId string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"id": userId})

	if err != nil {
		return errors.Wrap(err, "[UserRepository] Error deleting user "+userId)
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
