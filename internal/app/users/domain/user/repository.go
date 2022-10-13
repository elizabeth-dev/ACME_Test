package user

import (
	"context"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
)

type Repository interface {
	AddUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, userId string) (*User, error)
	GetUsers(
		ctx context.Context,
		filters []query_utils.Filter,
		sort []query_utils.Sort,
		pagination query_utils.Pagination,
	) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	RemoveUser(ctx context.Context, userId string) error
}
