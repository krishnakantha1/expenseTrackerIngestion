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

func HandleGetMonthlyData(db *mongo.Client, w http.ResponseWriter, r *http.Request) {
	req := t.MonthGetRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		util.BadRequestResponse(w, "issue while reading the data provided in body", err.Error())
		return
	}

	startYear, startMonth, err := util.GetYearMonth(req.Date)

	if err != nil {
		util.BadRequestResponse(w, err.Error())
		return
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		util.BadRequestResponse(w, "issue with getting time zone")
		return
	}

	startDate := time.Date(startYear, startMonth, 1, 0, 0, 0, 0, loc)

	endDateE := startDate.AddDate(0, 1, 0)

	expenses, err := dw.GetExpenseDataByDateRange(db, startDate, endDateE)

	if err != nil {
		util.BadRequestResponse(w, "issue while reading data")
	}

	res := t.MonthGetResponse{
		Error:    false,
		Count:    int64(len(*expenses)),
		Expenses: expenses,
	}

	json.NewEncoder(w).Encode(res)
}
