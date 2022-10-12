package query

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/domain/user"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
)

/*
The GetUsers query returns a list of users matching the given parameters.

The possible parameters include:
- Filters: a list of AND filters to apply to the query. If an OR filter is needed, multiple queries should be made.
- Sort: a list of fields to sort by, from most priority to least priority.
- Pagination: the pagination parameters to apply to the query.
*/
type GetUsers struct {
	Filters    []query_utils.Filter
	Sort       []query_utils.Sort
	Pagination query_utils.Pagination
}

type GetUsersHandler struct {
	userRepo user.Repository
}

func NewGetUsersHandler(userRepo user.Repository) *GetUsersHandler {
	if userRepo == nil {
		panic("[query/get_users] nil userRepo")
	}

	return &GetUsersHandler{userRepo}
}

func (h *GetUsersHandler) Handle(ctx context.Context, query GetUsers) ([]*User, error) {
	usersResult, err := h.userRepo.GetUsers(ctx, query.Filters, query.Sort, query.Pagination)

	if err != nil {
		return nil, err
	}

	var users []*User
	for _, u := range usersResult {
		users = append(
			users, &User{
				Id:        u.Id(),
				FirstName: u.FirstName(),
				LastName:  u.LastName(),
				Nickname:  u.Nickname(),
				Password:  u.Password(),
				Email:     u.Email(),
				Country:   u.Country(),
				CreatedAt: u.CreatedAt(),
				UpdatedAt: u.UpdatedAt(),
			},
		)
	}

	return users, nil
}
