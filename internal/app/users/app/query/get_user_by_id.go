package query

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/pkg/errors"
)

/*
The GetUserById query returns a single user matching the provided Id.
*/

type IGetUserByIdHandler interface {
	Handle(ctx context.Context, userId string) (*User, error)
}

type GetUserByIdHandler struct {
	userRepo user.UserRepository
}

func NewGetUserByIdHandler(userRepo user.UserRepository) *GetUserByIdHandler {
	if userRepo == nil {
		panic("[query/get_user_by_id] nil userRepo")
	}

	return &GetUserByIdHandler{userRepo}
}

func (h *GetUserByIdHandler) Handle(ctx context.Context, userId string) (*User, error) {
	userResult, err := h.userRepo.GetUserById(ctx, userId)

	if err != nil {
		return nil, errors.Wrap(err, "[query/get_user_by_id] Error retrieving user from database")
	}

	return &User{
		Id:        userResult.Id(),
		FirstName: userResult.FirstName(),
		LastName:  userResult.LastName(),
		Nickname:  userResult.Nickname(),
		Email:     userResult.Email(),
		Password:  userResult.Password(),
		Country:   userResult.Country(),
		CreatedAt: userResult.CreatedAt(),
		UpdatedAt: userResult.UpdatedAt(),
	}, nil
}
