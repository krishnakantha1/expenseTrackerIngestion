package datawrapper

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	da "github.com/krishnakantha1/expenseTrackerIngestion/database/data-access"
	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
)

func UpsertExpenseMessages(db *mongo.Client, expenseMessages []*t.ExpenseMessage, id primitive.ObjectID) {
	argsUpsert := t.UpsertAllArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "user_expenses",
		},
		SingleTransaction: make([]map[string]bson.M, 0, len(expenseMessages)),
	}

	expsenseModels := make([]*t.ExpsenseModel, 0, len(expenseMessages))

	for _, em := range expenseMessages {
		expsenseModels = append(expsenseModels, &t.ExpsenseModel{
			UserID:          id,
			URI:             em.ExpsenseEntry.URI,
			Bank:            em.ExpsenseEntry.Bank,
			AmountEncrypted: em.ExpsenseEntry.EncryptedAmount,
			ExpenseDate:     time.Unix((int64(em.ExpsenseEntry.ExpenseDate) / 1000), ((int64(em.ExpsenseEntry.ExpenseDate) % 1000) * 1000_000)),
			UpdatedOn:       time.Now(),
			ExpenseType:     em.ExpsenseEntry.ExpenseType,
			ExpenseTag:      em.ExpsenseEntry.ExpenseTag,
		})
	}

	for i, em := range expsenseModels {
		argsUpsert.SingleTransaction = append(argsUpsert.SingleTransaction, make(map[string]bson.M))

		argsUpsert.SingleTransaction[i]["updateValues"] = bson.M{"$set": bson.M{"user_id": em.UserID, "uri": em.URI,
			"bank": em.Bank, "amount_encrypted": em.AmountEncrypted,
			"expense_date": em.ExpenseDate, "updated_on": em.UpdatedOn,
			"expense_type": em.ExpenseType, "tag": em.ExpenseTag}}

		argsUpsert.SingleTransaction[i]["filter"] = bson.M{"user_id": em.UserID, "uri": em.URI}
	}

	da.UpsertAll(&argsUpsert)

	rawMessages := make([]*t.RawMessage, len(expenseMessages))
	for _, em := range expenseMessages {
		rawMessages = append(rawMessages, &em.RawMessage)
	}
	InsertRawMessages(db, rawMessages, "VALID")
}

func InsertRawMessages(db *mongo.Client, rawMessages []*t.RawMessage, rawtype string) {
	argsInsert := t.InsertAllArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "raw",
		},
		SingleTransaction: make([]interface{}, 0, len(rawMessages)),
	}

	for _, rm := range rawMessages {
		argsInsert.SingleTransaction = append(argsInsert.SingleTransaction, bson.D{{Key: "sms", Value: rm.Raw}, {Key: "type", Value: rawtype}})
	}

	_, err := da.InsertAll(&argsInsert)
	if err != nil {
		log.Println("err : ", err.Error())
	}
}

func GetExpenseDataByDateRange(db *mongo.Client, startDate time.Time, endDateE time.Time, userId primitive.ObjectID) (*[]t.ExpenseResponse, error) {
	argsSelectAll := t.SelectAllArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "user_expenses",
		},
		Filter: bson.M{
			"user_id": userId,
			"expense_date": bson.M{
				"$gte": startDate,
				"$lt":  endDateE,
			},
		},
	}

	cursor, err := da.SelectAll(&argsSelectAll)
	defer func() {
		err := cursor.Close(context.Background())
		if err != nil {
			log.Println("error while closing db cursor", err)
		}
	}()

	if err != nil {
		return nil, err
	}

	expenses := make([]t.ExpenseResponse, 0, 100)

	for cursor.Next(context.Background()) {
		expense := t.ExpenseResponse{}
		cursor.Decode(&expense)
		expenses = append(expenses, expense)
	}

	return &expenses, nil
}

func InsertSpamRaw() {}
