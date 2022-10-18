package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/sirupsen/logrus"
)

/*
The RemoveUser command removes a user from our platform given its id.
*/

type IRemoveUserHandler interface {
	Handle(ctx context.Context, userId string) error
}

type RemoveUserHandler struct {
	userRepo user.UserRepository
}

const removeUserTag = "command/remove_user"

func NewRemoveUserHandler(userRepo user.UserRepository) *RemoveUserHandler {
	if userRepo == nil {
		panic("[command/remove_user] nil userRepo")
	}

	return &RemoveUserHandler{userRepo}
}

func (h *RemoveUserHandler) Handle(ctx context.Context, userId string) error {
	logrus.WithFields(
		logrus.Fields{
			"tag":    removeUserTag,
			"userId": userId,
		},
	).Debug("Removing user")

	_, err := h.userRepo.GetUserById(ctx, userId)

	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":    removeUserTag,
				"userId": userId,
			},
		).WithError(err).Debug("Attempting to remove a nonexistent user")

		return err
	}

	if err := h.userRepo.RemoveUser(ctx, userId); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"tag":    removeUserTag,
				"userId": userId,
			},
		).WithError(err).Error("Error removing user")

		return err
	}

	return nil
}
