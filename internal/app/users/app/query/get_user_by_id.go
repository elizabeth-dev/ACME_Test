package query

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
)

/*
The GetUserById query returns a single user matching the provided Id.
*/

type GetUserByIdHandler struct {
	userRepo user.Repository
}

func NewGetUserByIdHandler(userRepo user.Repository) GetUserByIdHandler {
	if userRepo == nil {
		panic("[query/get_user_by_id] nil userRepo")
	}

	return GetUserByIdHandler{userRepo}
}

func (h *GetUserByIdHandler) Handle(ctx context.Context, userId string) (*User, error) {
	userResult, err := h.userRepo.GetUserById(ctx, userId)

	if err != nil {
		return nil, err
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