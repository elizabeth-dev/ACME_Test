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
	CreateUser command.ICreateUserHandler
	RemoveUser command.IRemoveUserHandler
	UpdateUser command.IUpdateUserHandler
}

type Queries struct {
	GetUsers    query.IGetUsersHandler
	GetUserById query.IGetUserByIdHandler
}
