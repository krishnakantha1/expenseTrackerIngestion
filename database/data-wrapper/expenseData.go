package datawrapper

import (
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

	argsInsert := t.InsertAllArgs{
		MongoArgs: t.MongoArgs{
			Client:     db,
			Database:   "ExpenceTracker",
			Collection: "raw",
		},
		SingleTransaction: make([]interface{}, 0, len(expenseMessages)),
	}

	for _, em := range expenseMessages {
		argsInsert.SingleTransaction = append(argsInsert.SingleTransaction, bson.D{{Key: "sms", Value: em.RawMessage.Raw}})
	}

	count, err := da.InsertAll(&argsInsert)
	if err != nil {
		log.Println("err : ", err.Error())
	}

	log.Println("--- count :", count)

}
