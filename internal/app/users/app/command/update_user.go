package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/pkg/errors"
)

/*
The UpdateUser command updates the given properties for a user in our platform, and leaves the rest of the properties untouched
*/
type UpdateUser struct {
	Id        string
	FirstName *string
	LastName  *string
	Nickname  *string
	Password  *string
	Email     *string
	Country   *string
}

type UpdateUserHandler struct {
	userRepo user.Repository
}

func NewUpdateUserHandler(userRepo user.Repository) *UpdateUserHandler {
	if userRepo == nil {
		panic("[command/update_user] nil userRepo")
	}

	return &UpdateUserHandler{userRepo}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUser) error {
	userToUpdate, err := h.userRepo.GetUserById(ctx, cmd.Id)
	if err != nil {
		return errors.Wrap(err, "[command/update_user] Error getting user "+cmd.Id+" from database")
	}

	if err := userToUpdate.Update(
		cmd.FirstName,
		cmd.LastName,
		cmd.Nickname,
		cmd.Password,
		cmd.Email,
		cmd.Country,
	); err != nil {
		return errors.Wrap(err, "[command/update_user] Error updating user "+cmd.Id)
	}

	if err := h.userRepo.UpdateUser(ctx, userToUpdate); err != nil {
		return errors.Wrap(err, "[command/update_user] Error updating user "+cmd.Id+" in database")
	}

	return nil
}
