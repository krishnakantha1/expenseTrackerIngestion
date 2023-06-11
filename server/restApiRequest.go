package server

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	h "github.com/krishnakantha1/expenseTrackerIngestion/api-handler"
	util "github.com/krishnakantha1/expenseTrackerIngestion/util"
)

type ApiHandlerFunc func(*mongo.Client, http.ResponseWriter, *http.Request)

func (s *Server) handleRestApiRequest() {
	//creating users end points
	http.HandleFunc("/api/login", s.handleRequestWithVerb(h.HandleLogin, http.MethodPost))
	http.HandleFunc("/api/register", s.handleRequestWithVerb(h.HandleRegister, http.MethodPost))

	//data ingestion end points

	//data providing end points
}

// handle function wrapper
func (s *Server) handleRequestWithVerb(f ApiHandlerFunc, verb string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method == verb {
			f(s.dbClient, w, r)
		} else {
			util.ForbidenResponse(w, "Method not supported.")
		}
	}
}
