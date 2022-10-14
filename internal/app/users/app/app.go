package app

import (
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/app/command"
	"github.com/elizabeth-dev/FACEIT_Test/internal/app/users/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateUser command.CreateUserHandler
	RemoveUser command.RemoveUserHandler
	UpdateUser command.UpdateUserHandler
}

type Queries struct {
	GetUsers query.GetUsersHandler
}
