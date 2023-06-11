package types

import "time"

type UserModel struct {
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedOn time.Time `bson:"createdon"`
	AesTest   string    `bson:"aestest"`
}
