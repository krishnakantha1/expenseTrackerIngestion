package types

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	AesTest  string `json:"aesTestString"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Username string `json:"username"`
	AesTest  string `json:"aesTestString"`
	JWT      string `json:"jwt"`
}

type UserData struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Username string             `json:"username" bson:"username"`
	AesTest  string             `json:"aesTestString" bson:"aestest"`
}

type JWTClaims struct {
	ID       primitive.ObjectID `json:"_id"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Username string             `json:"username"`

	jwt.RegisteredClaims
}
