package datawrapper

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	da "github.com/krishnakantha1/expenseTrackerIngestion/database/data-access"
	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
)

func GetUserDetail(db *mongo.Client, email string) *t.UserData {
	userDataAddr := new(t.UserData)

	args := t.SelectOneArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "users",
		},
		Filter: bson.M{"email": email},
		Opt: options.FindOne().SetProjection(bson.D{{Key: "email", Value: 1},
			{Key: "password", Value: 1},
			{Key: "username", Value: 1},
			{Key: "aestest", Value: 1}}),
		DecodeInto: userDataAddr,
	}

	da.SelectOne(&args)

	return userDataAddr
}

func SaveUserDetails(db *mongo.Client, user *t.UserModel) (primitive.ObjectID, error) {
	args := t.InsertOneArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "users",
		},
		Data: user,
	}

	ID, err := da.InsertOne(&args)

	return ID, err
}
