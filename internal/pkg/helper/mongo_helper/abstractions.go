package mongo_helper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) Collection
	Client() *mongo.Client
}

type Collection interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (cur Cursor, err error)
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type Cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
}

type MongoDatabase struct {
	Db *mongo.Database
}

func (md *MongoDatabase) Collection(colName string, opts ...*options.CollectionOptions) Collection {
	collection := md.Db.Collection(colName, opts...)
	return &MongoCollection{col: collection}
}

func (md *MongoDatabase) Client() *mongo.Client {
	return md.Db.Client()
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

func (sr *MongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

type MongoCursor struct {
	sr *mongo.Cursor
}

func (sr *MongoCursor) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (sr *MongoCursor) Next(ctx context.Context) bool {
	return sr.sr.Next(ctx)
}

type MongoCollection struct {
	col *mongo.Collection
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cur, err := mc.col.Find(ctx, filter, opts...)
	return &MongoCursor{sr: cur}, err
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := mc.col.FindOne(ctx, filter)
	return &MongoSingleResult{sr: singleResult}
}

func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.col.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *MongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.col.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *MongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (
	interface{}, error,
) {
	res, err := mc.col.UpdateOne(ctx, filter, update, opts...)
	return res, err
}
