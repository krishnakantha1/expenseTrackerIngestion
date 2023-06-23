package dataaccess

import (
	"context"

	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Select only one record from mongo db database
func SelectOne(args *t.SelectOneArgs) error {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	err := collection.FindOne(context.Background(), args.Filter, args.Opt).Decode(args.DecodeInto)

	return err
}

// Insert one record into mongo db database
func InsertOne(args *t.InsertOneArgs) (primitive.ObjectID, error) {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	result, err := collection.InsertOne(context.Background(), args.Data)

	if err != nil {
		return primitive.NewObjectID(), err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

// Upsert many records into mongo db database
func UpsertAll(args *t.UpsertAllArgs) (int64, error) {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	var bulk []mongo.WriteModel

	for i := 0; i < len(args.SingleTransaction); i++ {
		filter := args.SingleTransaction[i]["filter"]
		update := args.SingleTransaction[i]["updateValues"]

		updateModal := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)

		bulk = append(bulk, updateModal)
	}

	bulkWriteOptions := options.BulkWrite().SetOrdered(false)

	result, err := collection.BulkWrite(context.Background(), bulk, bulkWriteOptions)

	if err != nil {
		return 0, err
	}

	return result.InsertedCount + result.UpsertedCount, nil
}

func InsertAll(args *t.InsertAllArgs) (int64, error) {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	result, err := collection.InsertMany(context.Background(), args.SingleTransaction)

	return int64(len(result.InsertedIDs)), err
}

func SelectAll(args *t.SelectAllArgs) (*mongo.Cursor, error) {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	cursor, err := collection.Find(context.Background(), args.Filter)

	return cursor, err
}
