package user

import "context"

type Repository interface {
	AddUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	RemoveUser(ctx context.Context, userId string) error
}
