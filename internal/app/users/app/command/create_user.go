package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/google/uuid"
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

func NewCreateUserHandler(userRepo user.UserRepository) *CreateUserHandler {
	if userRepo == nil {
		panic("[command/create_user] nil userRepo")
	}

	return &CreateUserHandler{userRepo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) (string, error) {
	newId := uuid.NewString()

	newUser, err := user.CreateUser(
		newId,
		cmd.FirstName,
		cmd.LastName,
		cmd.Nickname,
		cmd.Password,
		cmd.Email,
		cmd.Country,
	)

	if err != nil {
		return "", err
	}

	if err := h.userRepo.AddUser(ctx, newUser); err != nil {
		return "", err
	}

	return newId, nil
}
