package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/sirupsen/logrus"
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

type IUpdateUserHandler interface {
	Handle(ctx context.Context, cmd UpdateUser) error
}

type UpdateUserHandler struct {
	userRepo user.UserRepository
}

const updateUserTag = "command/update_user"

func NewUpdateUserHandler(userRepo user.UserRepository) *UpdateUserHandler {
	if userRepo == nil {
		panic("[command/update_user] nil userRepo")
	}

	return &UpdateUserHandler{userRepo}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUser) error {
	logrus.WithFields(
		logrus.Fields{
			"tag": updateUserTag,
			"cmd": cmd,
		},
	).Debug("Updating user")

	userToUpdate, err := h.userRepo.GetUserById(ctx, cmd.Id)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag": updateUserTag,
				"cmd": cmd,
			},
		).WithError(err).Error("Error getting user to update")

		return err
	}

	if err := userToUpdate.Update(
		cmd.FirstName, cmd.LastName, cmd.Nickname, cmd.Password, cmd.Email, cmd.Country,
	); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":          updateUserTag,
				"cmd":          cmd,
				"userToUpdate": userToUpdate,
			},
		).WithError(err).Error("Error updating user")

		return err
	}

	if err := h.userRepo.UpdateUser(ctx, userToUpdate); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":          updateUserTag,
				"cmd":          cmd,
				"userToUpdate": userToUpdate,
			},
		).WithError(err).Error("Error updating user")

		return err
	}

	return nil
}
