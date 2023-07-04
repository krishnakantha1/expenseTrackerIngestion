package server

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"

	wsh "github.com/krishnakantha1/expenseTrackerIngestion/websocket-handler"
)

type WSHandlerFunc func(*mongo.Client, *websocket.Conn)

func (s *Server) handleWebsocketRequest() {

	http.HandleFunc("/ws/ingestExpenses", func(w http.ResponseWriter, req *http.Request) {
		ser := websocket.Server{
			Handler: websocket.Handler(s.handleRequest(wsh.IngestExpenseReadLoop)),
		}
		ser.ServeHTTP(w, req)
	})

	http.HandleFunc("/ws/ingestSpams", func(w http.ResponseWriter, req *http.Request) {
		ser := websocket.Server{
			Handler: websocket.Handler(s.handleRequest(wsh.IngestSpamMsgReadLoop)),
		}
		ser.ServeHTTP(w, req)
	})

}

func (s *Server) handleRequest(f WSHandlerFunc) websocket.Handler {

	return func(ws *websocket.Conn) {

		fmt.Println("new incomming connection from client:", ws.RemoteAddr())

		f(s.dbClient, ws)
	}
}
