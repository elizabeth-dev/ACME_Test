package user

import (
	pkgErrors "github.com/elizabeth-dev/ACME_Test/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

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
	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"getters": {
			"should return the correct values": testGetters,
		},
		"create user": {
			"create simple user":                    testCreateUser,
			"create user with each field empty":     testCreateUserWithFieldsEmpty,
			"create user with several fields empty": testCreateUserWithSeveralFieldsEmpty,
			"create user with hash fail":            testCreateUserWithHashFail,
		},
		"update user": {
			"update user":                           testUpdateUser,
			"update user with each field empty":     testUpdateUserWithEmptyFields,
			"update user with several fields empty": testUpdateUserWithSeveralEmptyFields,
			"update user with hash error":           testUpdateUserWithHashError,
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

	// Set the stubbed functions back to their original values so they don't affect other tests.
	nowFunc = time.Now
	hashFunc = bcrypt.GenerateFromPassword
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

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "id", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, "", lastName, nickname, password, email, country)

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "first_name", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, "", nickname, password, email, country)

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "last_name", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, "", password, email, country)

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "nickname", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, "", email, country)

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "password", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, password, "", country)

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "email", Value: ""}, err)
	assert.Nil(t, got)

	got, err = CreateUser(id, firstName, lastName, nickname, password, email, "")

	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "country", Value: ""}, err)
	assert.Nil(t, got)
}

func testCreateUserWithSeveralFieldsEmpty(t *testing.T) {
	lastName := "Doe"
	nickname := "john-123"
	password := "password"
	email := "me@john.com"
	country := "US"

	now := time.Now()
	setNow(now)

	hashedPassword := []byte("password")
	setHash(hashedPassword, nil)

	got, err := CreateUser("", "", lastName, nickname, password, email, country)

	assert.IsType(t, &pkgErrors.MultipleInvalidFields{}, err)
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

	hashErr := errors.New("hash fail")
	setHash(nil, hashErr)

	got, err := CreateUser(id, firstName, lastName, nickname, password, email, country)

	assert.Equal(
		t, &pkgErrors.Unknown{
			Tag:   domain,
			Cause: hashErr,
		}, err,
	)
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
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "first_name", Value: empty}, err)
	assert.Equal(t, User1, user)

	err = user.Update(&firstName, &empty, &nickname, &password, &email, &country)
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "last_name", Value: empty}, err)
	assert.Equal(t, User1, user)

	err = user.Update(&firstName, &lastName, &empty, &password, &email, &country)
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "nickname", Value: empty}, err)
	assert.Equal(t, User1, user)

	err = user.Update(&firstName, &lastName, &nickname, &empty, &email, &country)
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "password", Value: empty}, err)
	assert.Equal(t, User1, user)

	err = user.Update(&firstName, &lastName, &nickname, &password, &empty, &country)
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "email", Value: empty}, err)
	assert.Equal(t, User1, user)

	err = user.Update(&firstName, &lastName, &nickname, &password, &email, &empty)
	assert.Equal(t, &pkgErrors.InvalidField{Domain: "User", Field: "country", Value: empty}, err)
	assert.Equal(t, User1, user)
}

func testUpdateUserWithSeveralEmptyFields(t *testing.T) {
	user := User1

	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"

	empty := ""

	hashedPassword := []byte("hashed")
	setHash(hashedPassword, nil)

	err := user.Update(&empty, &empty, &nickname, &password, &email, &country)
	assert.IsType(t, &pkgErrors.MultipleInvalidFields{}, err)
	assert.Equal(t, User1, user)
}

func testUpdateUserWithHashError(t *testing.T) {
	user := User1

	firstName := "updated"
	lastName := "updated"
	nickname := "updated"
	password := "updated"
	email := "updated"
	country := "updated"

	hashErr := errors.New("hash fail")
	setHash(nil, hashErr)

	err := user.Update(&firstName, &lastName, &nickname, &password, &email, &country)

	assert.Equal(
		t, &pkgErrors.Unknown{
			Tag:   domain,
			Cause: hashErr,
		}, err,
	)
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
