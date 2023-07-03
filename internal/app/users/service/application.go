package service

import (
	"context"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/adapter"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/command"
	"github.com/elizabeth-dev/ACME_Test/internal/app/users/app/query"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/helper/mongo_helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"os"
)

func NewApplication(ctx context.Context) (app.Application, map[string]func(ctx context.Context) error) {
	dbClient := setupMongo(ctx)
	userRepo := adapter.NewUserRepository(dbClient)

	dependencies := map[string]func(ctx context.Context) error{
		"mongodb": func(ctx context.Context) error {
			return dbClient.Client().Ping(ctx, nil)
		},
	}

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(&userRepo),
			UpdateUser: command.NewUpdateUserHandler(&userRepo),
			RemoveUser: command.NewRemoveUserHandler(&userRepo),
		},
		Queries: app.Queries{
			GetUsers:    query.NewGetUsersHandler(&userRepo),
			GetUserById: query.NewGetUserByIdHandler(&userRepo),
		},
	}, dependencies

}

func setupMongo(ctx context.Context) mongo_helper.Database {
	mongoUri := os.Getenv("MONGODB_URI")

	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		panic(err)
	}

	_mongoUri, err := url.Parse(mongoUri)

	var db string

	// Check if the database is embedded in the MongoDB connection URI.
	if err == nil {
		db = _mongoUri.Path[1:]
	}

	// If the database is not embedded in the URI, then look for it in environment variables.
	if db == "" {
		db = os.Getenv("MONGODB_DB")
	}

	// If no database is found, then exit.
	if db == "" {
		log.Fatal("You must set your db on the MongoDB URI or in the 'MONGODB_DB' environmental variable.")
	}

	return &mongo_helper.MongoDatabase{Db: client.Database(db)}
}
