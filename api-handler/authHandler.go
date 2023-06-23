package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	dw "github.com/krishnakantha1/expenseTrackerIngestion/database/data-wrapper"
	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
	util "github.com/krishnakantha1/expenseTrackerIngestion/util"
)

// http handler for servicing user login request
func HandleLogin(db *mongo.Client, w http.ResponseWriter, r *http.Request) {
	//get Data from request body
	reqBody := t.LoginRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	if err != nil {
		util.BadRequestResponse(w, "issue while reading the data provided in body.")
		return
	}

	//get user details from database
	userData := dw.GetUserDetail(db, reqBody.Email)

	//check passwrod against hash
	isValidPassword := util.CheckHashedPassword(reqBody.Password, userData.Password)

	if !isValidPassword {
		util.BadRequestResponse(w, "username or password dosent exist.")
		return
	}

	//encode userData into jwt
	jwt, err := util.EncodeJWT(userData)

	if err != nil {
		util.BadRequestResponse(w, "Issue with JWT creation.")
		return
	}

	resp := t.AuthResponse{
		Username: userData.Username,
		AesTest:  userData.AesTest,
		JWT:      jwt,
	}
	json.NewEncoder(w).Encode(resp)

}

// http handler for servicing user register request
func HandleRegister(db *mongo.Client, w http.ResponseWriter, r *http.Request) {
	reqBody := t.RegisterRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	if err != nil {
		util.BadRequestResponse(w, "issue while reading the data provided in body.")
		return
	}

	hashedPassword, err := util.HashPassword(reqBody.Password)

	if err != nil {
		util.BadRequestResponse(w, "server error.(hashing password)")
		return
	}

	userModel := t.UserModel{
		Username:  reqBody.Username,
		Email:     reqBody.Email,
		Password:  hashedPassword,
		CreatedOn: time.Now(),
		AesTest:   reqBody.AesTest,
	}

	ID, err := dw.SaveUserDetails(db, &userModel)

	if err != nil {
		util.BadRequestResponse(w, err.Error())
		return
	}

	userData := t.UserData{
		ID:       ID,
		Email:    reqBody.Email,
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	jwt, err := util.EncodeJWT(&userData)

	if err != nil {
		util.BadRequestResponse(w, "Issue with JWT creation.")
		return
	}

	resp := t.AuthResponse{
		Username: reqBody.Username,
		AesTest:  reqBody.AesTest,
		JWT:      jwt,
	}
	json.NewEncoder(w).Encode(resp)
}
