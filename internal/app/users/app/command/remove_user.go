package command

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/pkg/errors"
)

/*
The RemoveUser command removes a user from our platform given its id.
*/
type RemoveUser struct {
	Id string
}

type RemoveUserHandler struct {
	userRepo user.Repository
}

func NewRemoveUserHandler(userRepo user.Repository) *RemoveUserHandler {
	if userRepo == nil {
		panic("[command/remove_user] nil userRepo")
	}

	return &RemoveUserHandler{userRepo}
}

func (h *RemoveUserHandler) Handle(ctx context.Context, cmd RemoveUser) error {
	if err := h.userRepo.RemoveUser(ctx, cmd.Id); err != nil {
		return errors.Wrap(err, "[command/remove_user] Error removing user "+cmd.Id+" from database")
	}

	return nil
}
