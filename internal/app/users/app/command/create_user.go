package command

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/domain/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

/*
The CreateUser command registers a new user into our platform and returns the generated id.
*/
type CreateUser struct {
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
}

type ICreateUserHandler interface {
	Handle(ctx context.Context, cmd CreateUser) (string, error)
}

type CreateUserHandler struct {
	userRepo user.UserRepository
}

const createUserTag = "command/create_user"

func NewCreateUserHandler(userRepo user.UserRepository) *CreateUserHandler {
	if userRepo == nil {
		panic("[command/create_user] nil userRepo")
	}

	return &CreateUserHandler{userRepo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) (string, error) {
	newId := uuid.NewString()

	logrus.WithFields(
		logrus.Fields{
			"tag":   createUserTag,
			"newId": newId,
			"cmd":   cmd,
		},
	).Debug("Creating user")

	newUser, err := user.CreateUser(
		newId, cmd.FirstName, cmd.LastName, cmd.Nickname, cmd.Password, cmd.Email, cmd.Country,
	)

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":   createUserTag,
				"newId": newId,
				"cmd":   cmd,
			},
		).WithError(err).Error("Error creating user")
		return "", err
	}

	if err := h.userRepo.AddUser(ctx, newUser); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":   createUserTag,
				"newId": newId,
				"cmd":   cmd,
			},
		).WithError(err).Error("Error calling repo AddUser")

		return "", err
	}

	return newId, nil
}
