package mongodb

import (
	"context"

	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SelectOne(args *t.SelectOneArgs) error {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	err := collection.FindOne(context.Background(), args.Filter, args.Opt).Decode(args.DecodeInto)

	return err
}

func InsertOne(args *t.InsertOneArgs) (primitive.ObjectID, error) {
	collection := args.Client.Database(args.Database).Collection(args.Collection)

	result, err := collection.InsertOne(context.Background(), args.Data)

	if err != nil {
		return primitive.NewObjectID(), err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

// func Select(client *mongo.Client, database string, collection string) {
// 	Collection := client.Database(database).Collection(collection)

// 	cur, err := Collection.Find(context.Background(), bson.D{{Key: "user_id", Value: 1234}})
// 	if err != nil {
// 		log.Fatal("Error in mongoDB/mongoDataAccess.go SelectOne : ", err)
// 	}
// 	defer cur.Close(context.Background())

// 	for cur.Next(context.Background()) {
// 		data := struct {
// 			Uid  int    `bson:"user_id"`
// 			Bank string `bson:"bank"`
// 		}{}

// 		cur.Decode(&data)

// 		fmt.Println(data.Uid, data.Bank)
// 	}
// }
