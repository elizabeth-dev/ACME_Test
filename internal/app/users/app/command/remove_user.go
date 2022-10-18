package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/pkg/errors"
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

func NewRemoveUserHandler(userRepo user.UserRepository) *RemoveUserHandler {
	if userRepo == nil {
		panic("[command/remove_user] nil userRepo")
	}

	return &RemoveUserHandler{userRepo}
}

func (h *RemoveUserHandler) Handle(ctx context.Context, userId string) error {
	_, err := h.userRepo.GetUserById(ctx, userId)

	if err != nil {
		return errors.Wrap(err, "[command/remove_user] Error retrieving user "+userId+" from database")
	}

	if err := h.userRepo.RemoveUser(ctx, userId); err != nil {
		return errors.Wrap(err, "[command/remove_user] Error removing user "+userId+" from database")
	}

	return nil
}
