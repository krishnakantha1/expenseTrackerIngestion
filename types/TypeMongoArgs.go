package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoArgs struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

type SelectOneArgs struct {
	MongoArgs
	Filter     bson.M
	Opt        *options.FindOneOptions
	DecodeInto any
}

type InsertOneArgs struct {
	MongoArgs
	Data any
}

type UpsertAllArgs struct {
	MongoArgs
	SingleTransaction []map[string]bson.M
}

type InsertAllArgs struct {
	MongoArgs
	SingleTransaction []interface{}
}

type SelectAllArgs struct {
	MongoArgs
	Filter bson.M
}
