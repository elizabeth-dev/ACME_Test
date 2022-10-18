package user

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var nowFunc = time.Now
var hashFunc = bcrypt.GenerateFromPassword

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

func (u *User) Update(
	firstName *string, lastName *string, nickname *string, password *string, email *string, country *string,
) error {
	/* Validation */
	var errorAcc []error
	if firstName != nil {
		if *firstName == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty first name"))
		}
	}

	if lastName != nil {
		if *lastName == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty last name"))
		}
	}

	if nickname != nil {
		if *nickname == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty nickname"))
		}
	}

	if password != nil {
		if *password == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty password"))
		}
	}

	if email != nil {
		if *email == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty email"))
		}
	}

	if country != nil {
		if *country == "" {
			errorAcc = append(errorAcc, errors.New("[User] Empty country"))
		}
	}

	if len(errorAcc) == 1 {
		return errorAcc[0]
	}

	if len(errorAcc) > 1 {
		return errors.New("[User] Multiple errors")
	}

	/* Update */

	if firstName != nil {
		u.firstName = *firstName
		u.updatedAt = nowFunc()
	}

	if lastName != nil {
		u.lastName = *lastName
		u.updatedAt = nowFunc()
	}

	if nickname != nil {
		u.nickname = *nickname
		u.updatedAt = nowFunc()
	}

	if password != nil {
		hashedPassword, err := hashPassword(*password)

		if err != nil {
			return errors.Wrap(err, "[User] Error hashing password")
		}

		u.password = hashedPassword
		u.updatedAt = nowFunc()
	}

	if email != nil {
		u.email = *email
		u.updatedAt = nowFunc()
	}

	if country != nil {
		u.country = *country
		u.updatedAt = nowFunc()
	}

	return nil
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
	var errorAcc []error

	if id == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty id"))
	}

	if firstName == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty first name"))
	}

	if lastName == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty last name"))
	}

	if nickname == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty nickname"))
	}

	if password == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty password"))
	}

	if email == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty email"))
	}

	if country == "" {
		errorAcc = append(errorAcc, errors.New("[User] Empty country"))
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		errorAcc = append(errorAcc, errors.Wrap(err, "[User] Error hashing password"))
	}

	if len(errorAcc) == 1 {
		return nil, errorAcc[0]
	}

	if len(errorAcc) > 1 {
		return nil, errors.New("[User] Multiple errors")
	}

	now := nowFunc()

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		nickname:  nickname,
		password:  hashedPassword,
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

/*
hashPassword is a helper method that takes a password and hashes it using the bcrypt hash function.

Now, here I've been making some research, as the OWASP foundation guidelines recommend using Argon2id, a newer hash function, but its strength compared to the bcrypt function seems to be debated under specific circumstances. Argon2id seems to be weaker to GPU attacks, but it's stronger than bcrypt against FPGA attacks. So... For now I think I'll stick with bcrypt, as it's still considered a strong hash function, has been field-tested for a longer time, and is also on the OWASP guidelines.

I've run a simple benchmark on bcrypt cost values. On my computer 13 rounds take ~600ms, while 14 rounds take ~1200ms. So I'm using 14 rounds, as it's closer to the general rule of 1 second.
*/
func hashPassword(password string) (string, error) {
	hashedPassword, err := hashFunc([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
