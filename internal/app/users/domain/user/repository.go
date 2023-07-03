package user

import (
	"context"
	"fmt"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils"
)

type NotFoundError struct {
	Id string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("User with id %s not found", e.Id)
}

/* UserRepository
Disclaimer: this should be called just "Repository". But doing so would mess up the mock generation with Mockery.
*/

type UserRepository interface {
	AddUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, userId string) (*User, error)
	GetUsers(
		ctx context.Context, filters []query_utils.Filter, sort []query_utils.Sort, pagination query_utils.Pagination,
	) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	RemoveUser(ctx context.Context, userId string) error
}
