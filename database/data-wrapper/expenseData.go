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
	args := t.UpsertAllArgs{
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
			ExpenseDate:     time.Unix(0, em.ExpsenseEntry.ExpenseDate),
			UpdatedOn:       time.Now(),
			ExpenseType:     em.ExpsenseEntry.ExpenseType,
			ExpenseTag:      em.ExpsenseEntry.ExpenseTag,
		})
	}

	for i, em := range expsenseModels {
		args.SingleTransaction = append(args.SingleTransaction, make(map[string]bson.M))

		args.SingleTransaction[i]["updateValues"] = bson.M{"$set": bson.M{"user_id": em.UserID, "uri": em.URI,
			"bank": em.Bank, "amount_encrypted": em.AmountEncrypted,
			"expense_date": em.ExpenseDate, "updated_on": em.UpdatedOn,
			"expense_type": em.ExpenseType, "tag": em.ExpenseTag}}

		args.SingleTransaction[i]["filter"] = bson.M{"user_id": em.UserID, "uri": em.URI}
	}

	count, err := da.UpsertAll(&args)

	log.Println("saved : ", count, "err : ", err)
}
