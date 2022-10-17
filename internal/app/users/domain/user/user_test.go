package user

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setNow(t time.Time) {
	nowFunc = func() time.Time {
		return t
	}
}

func setHash(h []byte, err error) {
	hashFunc = func(password []byte, cost int) ([]byte, error) {
		return h, err
	}
}

func TestUser(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"getters": {
			"should return the correct values": testGetters,
		},
		"create user": {
			"create simple user":                testCreateUser,
			"create user with each field empty": testCreateUserWithFieldsEmpty,
			"create user with hash fail":        testCreateUserWithHashFail,
		},
		"update user": {
			"update user":                       testUpdateUser,
			"update user with each field empty": testUpdateUserWithEmptyFields,
			"update user with hash error":       testUpdateUserWithHashError,
		},
		"unmarshal user": {
			"unmarshal user": testUnmarshalUser,
		},
	} {
		testGroup := testGroup
		t.Run(
			name, func(t *testing.T) {
				for name, test := range testGroup {
					test := test
					t.Run(
						name, func(t *testing.T) {
							test(t)
						},
					)
				}
			},
		)
	}
}

func testGetters(t *testing.T) {
	user := User1

	assert.Equal(t, user.id, user.Id())
	assert.Equal(t, user.firstName, user.FirstName())
	assert.Equal(t, user.lastName, user.LastName())
	assert.Equal(t, user.nickname, user.Nickname())
	assert.Equal(t, user.password, user.Password())
	assert.Equal(t, user.email, user.Email())
	assert.Equal(t, user.country, user.Country())
	assert.Equal(t, user.createdAt, user.CreatedAt())
	assert.Equal(t, user.updatedAt, user.UpdatedAt())
}

func testCreateUser(t *testing.T) {
	id := uuid.NewString()
	firstName := "John"
	lastName := "Doe"
	nickname := "john-123"
	password := "password"
	email := "me@john.com"
	country := "US"

	now := time.Now()
	setNow(now)

	hashedPassword := []byte("password")
	setHash(hashedPassword, nil)

	got, err := CreateUser(id, firstName, lastName, nickname, password, email, country)

	assert.NoError(t, err)

	expected := &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		nickname:  nickname,
		password:  string(hashedPassword),
		email:     email,
		country:   country,
		createdAt: now,
		updatedAt: now,
	}

	assert.Equal(t, expected, got)
}

func testCreateUserWithFieldsEmpty(t *testing.T) {
	id := uuid.NewString()
	firstName := "John"
	lastName := "Doe"
	nickname := "john-123"
	password := "password"
	email := "me@john.com"
	country := "US"

	now := time.Now()
	setNow(now)

	hashedPassword := []byte("password")
	setHash(hashedPassword, nil)

	got, err := CreateUser("", firstName, lastName, nickname, password, email, country)

	assert.EqualError(t, err, "[User] Empty id")
	assert.Nil(t, got)

	got, err = CreateUser(id, "", lastName, nickname, password, email, country)

	assert.EqualError(t, err, "[User] Empty first name")
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, "", nickname, password, email, country)

	assert.EqualError(t, err, "[User] Empty last name")
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, "", password, email, country)

	assert.EqualError(t, err, "[User] Empty nickname")
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, "", email, country)

	assert.EqualError(t, err, "[User] Empty password")
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, password, "", country)

	assert.EqualError(t, err, "[User] Empty email")
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, password, email, "")

	assert.EqualError(t, err, "[User] Empty country")
	assert.Nil(t, got)
}

func testCreateUserWithHashFail(t *testing.T) {
	id := uuid.NewString()
	firstName := "John"
	lastName := "Doe"
	nickname := "john-123"
	password := "password"
	email := "me@john.com"
	country := "US"

	now := time.Now()
	setNow(now)

	setHash(nil, errors.New("hash fail"))

	got, err := CreateUser(id, firstName, lastName, nickname, password, email, country)

	assert.EqualError(t, err, "[User] Error hashing password: hash fail")
	assert.Nil(t, got)
}

func testUpdateUser(t *testing.T) {
	user := User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"

	now := time.Now()
	setNow(now)

	hashedPassword := []byte("hashed")
	setHash(hashedPassword, nil)

	err := user.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	assert.NoError(t, err)
	assert.Equal(t, User1.id, user.id)
	assert.Equal(t, firstName, user.firstName)
	assert.Equal(t, lastName, user.lastName)
	assert.Equal(t, nickname, user.nickname)
	assert.Equal(t, string(hashedPassword), user.password)
	assert.Equal(t, email, user.email)
	assert.Equal(t, country, user.country)
	assert.Equal(t, User1.createdAt, user.createdAt)
	assert.Equal(t, now, user.updatedAt)
}

func testUpdateUserWithEmptyFields(t *testing.T) {
	user := User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"

	empty := ""

	hashedPassword := []byte("hashed")
	setHash(hashedPassword, nil)

	err := user.Update(&empty, &lastName, &nickname, &password, &email, &country)
	assert.EqualError(t, err, "[User] Empty first name")

	err = user.Update(&firstName, &empty, &nickname, &password, &email, &country)
	assert.EqualError(t, err, "[User] Empty last name")

	err = user.Update(&firstName, &lastName, &empty, &password, &email, &country)
	assert.EqualError(t, err, "[User] Empty nickname")

	err = user.Update(&firstName, &lastName, &nickname, &empty, &email, &country)
	assert.EqualError(t, err, "[User] Empty password")

	err = user.Update(&firstName, &lastName, &nickname, &password, &empty, &country)
	assert.EqualError(t, err, "[User] Empty email")

	err = user.Update(&firstName, &lastName, &nickname, &password, &email, &empty)
	assert.EqualError(t, err, "[User] Empty country")
}

func testUpdateUserWithHashError(t *testing.T) {
	user := User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"

	setHash(nil, errors.New("hash fail"))

	err := user.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	assert.EqualError(t, err, "[User] Error hashing password: hash fail")
}

func testUnmarshalUser(t *testing.T) {
	now := time.Now()

	id := "1234"
	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"
	createdAt := now
	updatedAt := now

	out := UnmarshalUserFromDB(id, firstName, lastName, nickname, password, email, country, createdAt, updatedAt)

	assert.Equal(t, id, out.id)
	assert.Equal(t, firstName, out.firstName)
	assert.Equal(t, lastName, out.lastName)
	assert.Equal(t, nickname, out.nickname)
	assert.Equal(t, password, out.password)
	assert.Equal(t, email, out.email)
	assert.Equal(t, country, out.country)
	assert.Equal(t, createdAt, out.createdAt)
	assert.Equal(t, updatedAt, out.updatedAt)
}
