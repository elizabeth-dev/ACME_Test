package user

import (
	"errors"
	"time"
)

/*
A User holds our domain model for a user entity.

This should reflect our business logic and needs. It should never be influenced
neither by the database requirements, nor by the API requirements.

It should also always keep valid state, this way we can avoid having to check
for unexpected behavior in the other layers, and improve scalability and
maintainability.
*/
type User struct {
	id        string
	firstName string
	lastName  string
	nickname  string
	password  string
	email     string
	country   string
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) Id() string {
	return u.id
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Country() string {
	return u.country
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

/*
CreateUser is the method we use to register new users into our platform.

The "register new users" part should be emphasized, as this should never be used
to unmarshal data from the database model, or to build a User struct from any
other source. These should be handled by its own "constructor" methods.

The reason for this is to preserve the rule of "keeping a valid state in the
domain layer". CreateUser holds the business logic required when a user signs up
in our platform, like validating the data against business rules, or setting
specific properties like createdAt.
*/
func CreateUser(
	id string, firstName string, lastName string, nickname string, password string, email string, country string,
) (*User, error) {
	if id == "" {
		return nil, errors.New("[User] Empty id")
	}

	if firstName == "" {
		return nil, errors.New("[User] Empty first name")
	}

	if lastName == "" {
		return nil, errors.New("[User] Empty last name")
	}

	if nickname == "" {
		return nil, errors.New("[User] Empty nickname")
	}

	if password == "" {
		return nil, errors.New("[User] Empty password")
	}

	if email == "" {
		return nil, errors.New("[User] Empty email")
	}

	if country == "" {
		return nil, errors.New("[User] Empty country")
	}

	now := time.Now()

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		nickname:  nickname,
		password:  password,
		email:     email,
		country:   country,
		createdAt: now,
		updatedAt: now,
	}, nil
}

/*
UnmarshalUserFromDB is the method we use to unmarshal data from the database model to the domain model.

As always in the domain layer, this process should be agnostic to the database
we're using, so we expose a method that takes the required parameters to build a
User struct.

Unlike the CreateUser method, this one does not apply any business validation
rules, as this will be used to parse data from the database only. Data that
should've been validated already before being stored in the database. Since the
database is the source of truth, we should trust the data we're receiving from
it.
*/
func UnmarshalUserFromDB(
	id string,
	firstName string,
	lastName string,
	nickname string,
	password string,
	email string,
	country string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		nickname:  nickname,
		password:  password,
		email:     email,
		country:   country,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
