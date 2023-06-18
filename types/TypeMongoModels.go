package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	USERS         = "users"
	USER_EXPENSES = "user_expenses"
	RAW_MESSAGES  = "raw_messages"
)

type UserModel struct {
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedOn time.Time `bson:"created_on"`
	UpdatedOn time.Time `bson:"updated_on"`
	AesTest   string    `bson:"aestest"`
}

type ExpsenseModel struct {
	UserID          primitive.ObjectID `bson:"user_id"`
	URI             string             `bson:"uri"`
	Bank            string             `bson:"bank"`
	AmountEncrypted string             `bson:"amount_encrypted"`
	ExpenseDate     string             `bson:"expense_date"`
	ExpenseType     string             `bson:"expense_type"`
	ExpenseTag      string             `bson:"tag"`
}

type RawMessageModel struct {
	Raw string `bson:"raw"`
}
