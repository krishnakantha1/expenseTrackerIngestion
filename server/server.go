package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
	//"github.com/rs/cors"
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

	err := godotenv.Load()
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

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if len(port) == 1 {
		port = ":8080"
	}
	// //cors
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"}, // Allow all origins
	// 	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
	// 	AllowCredentials: true,
	// })
	// //server
	// handler := websocket.Handler(s.HandleServer)

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		ser := websocket.Server{
			Handler: websocket.Handler(s.HandleServer),
		}
		ser.ServeHTTP(w, req)
	})
	http.ListenAndServe(port, nil)

}

func (s *Server) Select() {
	mdb.Select(s.dbClient, "ExpenceTracker", "user_expenses")
}
