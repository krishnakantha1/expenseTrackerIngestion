package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/websocket"
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
		log.Println("Error loading .env file", err)
	}

	//Get connectio to MongoDB
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	dbUrl := os.Getenv("DB_URL")
	dbOption := options.Client().ApplyURI(dbUrl)
	dbClient, err := mongo.Connect(context.Background(), dbOption)

	if err != nil {
		log.Fatal("Error while connecting to DB : ", err)
	}
	s.dbClient = dbClient

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if len(port) == 1 {
		port = ":8080"
	}

	log.Println("--server started on port", port)

	s.handleWebsocketRequest()
	s.handleRestApiRequest()

	if err = http.ListenAndServe(port, nil); err != nil {
		log.Fatal("unable to start server :", err)
	}
}
