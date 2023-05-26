package server

import (
	"context"
	"log"
	"net/http"

	//"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/websocket"

	mdb "github.com/krishnakantha1/expenseTrackerIngestion/database/mongoDB"
)

type Server struct {
	dbClient *mongo.Client
	conns    map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Init() {

	err := godotenv.Load("")
	if err != nil {
		//log.Fatal("Error loading .env file", err)
	}

	dbUrl := os.Getenv("DB_URL")

	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()

	dbClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal("Error while connecting to DB : ", err)
	}
	s.dbClient = dbClient

	//server
	http.Handle("/ws", websocket.Handler(s.HandleServer))
	http.ListenAndServe(":8080", nil)

}

func (s *Server) Select() {
	mdb.Select(s.dbClient, "ExpenceTracker", "user_expenses")
}
