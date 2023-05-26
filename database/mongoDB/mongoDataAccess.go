package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Select(client *mongo.Client, database string, collection string) {
	Collection := client.Database(database).Collection(collection)

	cur, err := Collection.Find(context.Background(), bson.D{{"user_id", 1234}})
	if err != nil {
		log.Fatal("Error in mongoDB/mongoDataAccess.go SelectOne : ", err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := struct {
			Uid  int    `bson:"user_id"`
			Bank string `bson:"bank"`
		}{}

		cur.Decode(&data)

		fmt.Println(data.Uid, data.Bank)
	}
}
