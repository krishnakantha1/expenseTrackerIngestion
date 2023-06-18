package util

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
)

func keyFunc(t *jwt.Token) (interface{}, error) {
	key := os.Getenv("PRIVATEKEY")

	if len(key) == 0 {
		log.Fatal("Secret key not found")
	}

	return []byte(key), nil
}

func EncodeJWT(userData *t.UserData) (string, error) {
	claims := t.JWTClaims{
		ID:       userData.ID,
		Username: userData.Username,
		Email:    userData.Email,
		Password: userData.Password,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := os.Getenv("PRIVATEKEY")

	if len(key) == 0 {
		log.Fatal("Secret key not found")
	}

	jwtString, err := token.SignedString([]byte(key))

	return jwtString, err
}

func DecodeJWT(jwtString string) (*t.UserData, error) {
	claims := t.JWTClaims{}

	log.Println(jwtString)

	token, err := jwt.ParseWithClaims(jwtString, claims, keyFunc)

	log.Println(claims)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("issue whith signature.")
		}

		return nil, errors.New("Server issue(DecodeJWT).")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	userData := t.UserData{
		ID:       claims.ID,
		Username: claims.Username,
		Email:    claims.Email,
		Password: claims.Password,
	}

	return &userData, nil
}
