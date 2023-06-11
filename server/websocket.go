package server

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func (s *Server) handleWebsocketRequest() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		ser := websocket.Server{
			Handler: websocket.Handler(s.HandleServer),
		}
		ser.ServeHTTP(w, req)
	})

}

func (s *Server) HandleServer(ws *websocket.Conn) {
	fmt.Println("new incomming connection from client:", ws.RemoteAddr())

	//s.conns[ws] = true

	s.ReadLoop(ws)
}

func (s *Server) ReadLoop(ws *websocket.Conn) {
	buf := make([]byte, 2048)
	defer ws.Close()

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("malformed message", err)
			return
		}
		msg := string(buf[:n])
		fmt.Println(msg)

		ws.Write([]byte("great! got your message"))
	}
}
